package handlers

import (
	"github.com/junjiexh/cowatching/internal/database"
)

type Handlers struct {
	Health *HealthHandler
	// Add more handlers here as needed
	// Users  *UsersHandler
	// Rooms  *RoomsHandler
	// Videos *VideosHandler
}

func New(db *database.Database) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(db),
		// Initialize other handlers here
	}
}
