package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/junjiexh/cowatching/internal/database"
	"github.com/junjiexh/cowatching/internal/database/db"
)

const (
	maxUploadSize = 500 << 20 // 500 MB
	uploadsDir    = "./uploads/videos"
)

type VideoHandler struct {
	db          *database.Database
	queries     *db.Queries
	uploadsPath string
}

type VideoResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

func NewVideoHandler(database *database.Database) *VideoHandler {
	// Ensure uploads directory exists
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create uploads directory: %v", err))
	}

	return &VideoHandler{
		db:          database,
		queries:     db.NewWithPool(database.Pool),
		uploadsPath: uploadsDir,
	}
}

// Upload handles video file uploads
func (h *VideoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to read video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "video/") {
		http.Error(w, "File must be a video", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)
	filePath := filepath.Join(h.uploadsPath, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save video", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	size, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "Failed to save video", http.StatusInternalServerError)
		return
	}

	// Get title from form or use filename
	title := r.FormValue("title")
	if title == "" {
		title = strings.TrimSuffix(header.Filename, ext)
	}

	// Save to database
	video, err := h.queries.CreateUploadedVideo(r.Context(), db.CreateUploadedVideoParams{
		Title:       title,
		Filename:    filename,
		ContentType: contentType,
		FileSize:    size,
	})
	if err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "Failed to save video metadata", http.StatusInternalServerError)
		return
	}

	// Create response
	response := VideoResponse{
		ID:          video.ID,
		Title:       video.Title,
		URL:         fmt.Sprintf("/api/v1/videos/stream/%d", video.ID),
		Size:        video.FileSize,
		ContentType: video.ContentType,
		UploadedAt:  video.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// List returns all uploaded videos
func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	videos, err := h.queries.ListUploadedVideos(r.Context())
	if err != nil {
		http.Error(w, "Failed to list videos", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	responses := make([]VideoResponse, 0, len(videos))
	for _, video := range videos {
		responses = append(responses, VideoResponse{
			ID:          video.ID,
			Title:       video.Title,
			URL:         fmt.Sprintf("/api/v1/videos/stream/%d", video.ID),
			Size:        video.FileSize,
			ContentType: video.ContentType,
			UploadedAt:  video.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// Stream serves a video file by ID
func (h *VideoHandler) Stream(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Video ID required", http.StatusBadRequest)
		return
	}

	videoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// Get video metadata from database
	video, err := h.queries.GetUploadedVideoByID(r.Context(), videoID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Video not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get video", http.StatusInternalServerError)
		}
		return
	}

	// Build file path
	filePath := filepath.Join(h.uploadsPath, video.Filename)

	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "Video file not found", http.StatusNotFound)
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open video", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set appropriate headers
	w.Header().Set("Content-Type", video.ContentType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	// Support range requests for video seeking
	http.ServeContent(w, r, video.Filename, info.ModTime(), file)
}

// Delete removes a video by ID
func (h *VideoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Video ID required", http.StatusBadRequest)
		return
	}

	videoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// Get video metadata from database
	video, err := h.queries.GetUploadedVideoByID(r.Context(), videoID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Video not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get video", http.StatusInternalServerError)
		}
		return
	}

	// Delete from database first
	if err := h.queries.DeleteUploadedVideo(r.Context(), videoID); err != nil {
		http.Error(w, "Failed to delete video", http.StatusInternalServerError)
		return
	}

	// Delete file from filesystem
	filePath := filepath.Join(h.uploadsPath, video.Filename)
	if err := os.Remove(filePath); err != nil {
		// Log error but don't fail the request since DB record is already deleted
		// In production, you might want to queue this for retry
		fmt.Printf("Warning: Failed to delete video file %s: %v\n", filePath, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Video deleted successfully",
	})
}
