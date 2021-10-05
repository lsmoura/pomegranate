package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"pomegranate/models"
	"pomegranate/themoviedb"

	bolt "go.etcd.io/bbolt"
)

type MovieEntry struct {
	Runtime  int32
	Released string
	ImdbId   string
	TmdbId   int32
	Year     int32
	Genres   []string
	Titles   []string
	Images   struct {
		Posters []string
	}
}

func (m *Manager) MovieSearch(tmdb themoviedb.Themoviedb, query string) ([]MovieEntry, error) {
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	res, err := tmdb.ReadMovies(query, 0)
	if err != nil {
		return nil, fmt.Errorf("tmdb.ReadMovies: %w", err)
	}

	var payload []MovieEntry

	for _, movie := range res.Results {
		m := MovieEntry{
			Titles:   []string{movie.Title},
			Released: movie.ReleaseDate,
			TmdbId:   movie.Id,
		}

		extraInfo, err := tmdb.ReadSingleMovie(fmt.Sprintf("%d", movie.Id))
		if err != nil {
			return nil, fmt.Errorf("tmdb.ReadSingleMovie (%d): %w", movie.Id, err)
		}

		m.ImdbId = extraInfo.ImdbId
		m.Runtime = extraInfo.Runtime

		payload = append(payload, m)
	}

	return payload, nil
}

func (m *Manager) MovieWithNzbID(id string) (models.Movie, error) {
	if m.DB.Database == nil {
		return models.Movie{}, fmt.Errorf("database was not initialized")
	}

	var resp models.Movie
	err := m.DB.Database.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(resp.Kind()))

		c := b.Cursor()

		for k, bytes := c.First(); k != nil; k, bytes = c.Next() {
			var m models.Movie
			if err := json.Unmarshal(bytes, &m); err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}

			for _, info := range m.NzbInfo {
				if info.ID == id {
					resp = m
					return nil
				}
			}
		}

		return nil
	})
	if err != nil {
		return models.Movie{}, err
	}

	return resp, nil
}

func (m *Manager) AllMovies() ([]models.Movie, error) {
	if m.DB.Database == nil {
		return nil, fmt.Errorf("database was not initialized")
	}

	var resp []models.Movie

	if err := m.Movies.FindAll(context.Background(), &resp); err != nil {
		return nil, fmt.Errorf("m.Movies.FindAll: %w", err)
	}

	return resp, nil
}

func (m *Manager) Movie(key string) (models.Movie, error) {
	if m.DB.Database == nil {
		return models.Movie{}, fmt.Errorf("database was not initialized")
	}

	var movie models.Movie
	if err := m.Movies.FindByID(context.Background(), &movie, key); err != nil {
		return models.Movie{}, fmt.Errorf("m.Movies.FindByID: %w", err)
	}

	return movie, nil
}
