package manager

import "pomegranate/database"

type Manager struct {
	database.DB

	Movies database.Store
}

func NewManager(db database.DB) (*Manager, error) {
	m := &Manager{
		DB:     db,
		Movies: database.NewStore(&db, &database.Movie{}),
	}

	return m, nil
}