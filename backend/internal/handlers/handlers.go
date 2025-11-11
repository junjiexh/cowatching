package handlers

import (
	"github.com/junjiexh/cowatching/internal/database"
)

type Handlers struct {
	Health *HealthHandler
	Video  *VideoHandler
	// Add more handlers here as needed
	// Users  *UsersHandler
	// Rooms  *RoomsHandler
}

func New(db *database.Database) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(db),
		Video:  NewVideoHandler(),
		// Initialize other handlers here
	}
}
