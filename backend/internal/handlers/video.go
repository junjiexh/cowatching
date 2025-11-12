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
	"github.com/junjiexh/cowatching/internal/s3"
)

const (
	maxUploadSize = 500 << 20 // 500 MB
	uploadsDir    = "./uploads/videos"
)

type VideoHandler struct {
	db          *database.Database
	queries     *db.Queries
	uploadsPath string
	s3Client    *s3.S3Client
}

type VideoResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

func NewVideoHandler(database *database.Database, s3Client *s3.S3Client) *VideoHandler {
	// Ensure uploads directory exists
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create uploads directory: %v", err))
	}

	return &VideoHandler{
		db:          database,
		queries:     db.NewWithPool(database.Pool),
		uploadsPath: uploadsDir,
		s3Client:    s3Client,
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

	// Generate unique filename/key for S3
	ext := filepath.Ext(header.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)

	// Get file size
	fileSize := header.Size

	// Get title from form or use filename
	title := r.FormValue("title")
	if title == "" {
		title = strings.TrimSuffix(header.Filename, ext)
	}

	// Upload to S3
	s3URL, err := h.s3Client.UploadVideo(r.Context(), filename, file, contentType, fileSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload video to S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Save metadata to database
	video, err := h.queries.CreateUploadedVideo(r.Context(), db.CreateUploadedVideoParams{
		Title:       title,
		Filename:    &filename,
		ContentType: contentType,
		FileSize:    fileSize,
		S3Key:       &filename,
		S3Url:       &s3URL,
	})
	if err != nil {
		// Try to clean up S3 file on database error
		if deleteErr := h.s3Client.DeleteVideo(r.Context(), filename); deleteErr != nil {
			fmt.Printf("Warning: Failed to clean up S3 file after DB error: %v\n", deleteErr)
		}
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

// Stream serves a video file by ID from S3
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

	// Check if video has S3 key
	if video.S3Key == nil || *video.S3Key == "" {
		http.Error(w, "Video not stored in S3", http.StatusNotFound)
		return
	}

	// Generate presigned URL with 1 hour expiration
	presignedURL, err := h.s3Client.GetVideoURL(r.Context(), *video.S3Key, 1*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate video URL", http.StatusInternalServerError)
		return
	}

	// Redirect to the presigned URL
	http.Redirect(w, r, presignedURL, http.StatusTemporaryRedirect)
}

// Delete removes a video by ID from database and S3
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

	// Delete file from S3
	if video.S3Key != nil && *video.S3Key != "" {
		if err := h.s3Client.DeleteVideo(r.Context(), *video.S3Key); err != nil {
			// Log error but don't fail the request since DB record is already deleted
			// In production, you might want to queue this for retry
			fmt.Printf("Warning: Failed to delete video from S3 %s: %v\n", *video.S3Key, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Video deleted successfully",
	})
}
