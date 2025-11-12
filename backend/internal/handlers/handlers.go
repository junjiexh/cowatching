package handlers

import (
	"github.com/junjiexh/cowatching/internal/database"
	"github.com/junjiexh/cowatching/internal/s3"
)

type Handlers struct {
	Health *HealthHandler
	Video  *VideoHandler
	// Add more handlers here as needed
	// Users  *UsersHandler
	// Rooms  *RoomsHandler
}

func New(db *database.Database, s3Client *s3.S3Client) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(db),
		Video:  NewVideoHandler(db, s3Client),
		// Initialize other handlers here
	}
}
