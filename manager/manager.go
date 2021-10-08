package manager

import (
	"pomegranate/database"
	"pomegranate/models"
)

type Manager struct {
	*database.DB

	Movies database.Store
}

func NewManager(db *database.DB) (*Manager, error) {
	m := &Manager{
		DB:     db,
		Movies: database.NewStore(db, &models.Movie{}),
	}

	return m, nil
}
