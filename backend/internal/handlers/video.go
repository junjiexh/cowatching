package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	maxUploadSize = 500 << 20 // 500 MB
	uploadsDir    = "./uploads/videos"
)

type VideoHandler struct {
	uploadsPath string
}

type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Filename    string    `json:"filename"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

func NewVideoHandler() *VideoHandler {
	// Ensure uploads directory exists
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create uploads directory: %v", err))
	}

	return &VideoHandler{
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
		http.Error(w, "Failed to save video", http.StatusInternalServerError)
		return
	}

	// Get title from form or use filename
	title := r.FormValue("title")
	if title == "" {
		title = strings.TrimSuffix(header.Filename, ext)
	}

	// Create response
	video := Video{
		ID:          fmt.Sprintf("%d", timestamp),
		Title:       title,
		Filename:    filename,
		URL:         fmt.Sprintf("/api/v1/videos/stream/%s", filename),
		Size:        size,
		ContentType: contentType,
		UploadedAt:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

// List returns all uploaded videos
func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(h.uploadsPath)
	if err != nil {
		http.Error(w, "Failed to read videos", http.StatusInternalServerError)
		return
	}

	videos := []Video{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// Extract timestamp from filename
		parts := strings.SplitN(file.Name(), "_", 2)
		var timestamp int64
		var title string
		if len(parts) == 2 {
			fmt.Sscanf(parts[0], "%d", &timestamp)
			ext := filepath.Ext(parts[1])
			title = strings.TrimSuffix(parts[1], ext)
		} else {
			timestamp = info.ModTime().Unix()
			ext := filepath.Ext(file.Name())
			title = strings.TrimSuffix(file.Name(), ext)
		}

		video := Video{
			ID:         fmt.Sprintf("%d", timestamp),
			Title:      title,
			Filename:   file.Name(),
			URL:        fmt.Sprintf("/api/v1/videos/stream/%s", file.Name()),
			Size:       info.Size(),
			UploadedAt: info.ModTime(),
		}

		videos = append(videos, video)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// Stream serves a video file
func (h *VideoHandler) Stream(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	// Security: prevent path traversal
	filename = filepath.Base(filename)
	filePath := filepath.Join(h.uploadsPath, filename)

	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
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
	w.Header().Set("Content-Type", "video/mp4") // Default to mp4, could be enhanced
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	// Support range requests for video seeking
	http.ServeContent(w, r, filename, info.ModTime(), file)
}

// Delete removes a video file
func (h *VideoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	// Security: prevent path traversal
	filename = filepath.Base(filename)
	filePath := filepath.Join(h.uploadsPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Failed to delete video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Video deleted successfully",
	})
}
