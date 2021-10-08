package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"pomegranate/models"
	"pomegranate/themoviedb"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type MovieEntry struct {
	Runtime  int32    `json:"runtime"`
	Released string   `json:"released"`
	ImdbId   string   `json:"imdb_id"`
	TmdbId   int32    `json:"tmdb_id"`
	Year     int32    `json:"year"`
	Genres   []string `json:"genres"`
	Titles   []string `json:"titles"`
	Images   struct {
		Posters []string `json:"posters"`
	} `json:"images"`
}

func (m *Manager) MovieSearch(tmdb themoviedb.Themoviedb, query string) ([]MovieEntry, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}

	res, err := tmdb.ReadMovies(query, 0)
	if err != nil {
		return nil, errors.Wrap(err, "tmdb.ReadMovies")
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
			return nil, errors.Wrapf(err, "tmdb.ReadSingleMovie (%d)", movie.Id)
		}

		m.ImdbId = extraInfo.ImdbId
		m.Runtime = extraInfo.Runtime

		payload = append(payload, m)
	}

	return payload, nil
}

func (m *Manager) MovieWithNzbID(id string) (models.Movie, error) {
	if m.DB.Database == nil {
		return models.Movie{}, errors.New("database was not initialized")
	}

	var resp models.Movie
	err := m.DB.Database.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(resp.Kind()))

		c := b.Cursor()

		for k, bytes := c.First(); k != nil; k, bytes = c.Next() {
			var m models.Movie
			if err := json.Unmarshal(bytes, &m); err != nil {
				return errors.Wrap(err, "json.Unmarshal")
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
		return models.Movie{}, errors.Wrap(err, "m.DB.Database.View")
	}

	return resp, nil
}

func (m *Manager) AllMovies() ([]models.Movie, error) {
	if m.DB.Database == nil {
		return nil, errors.New("database was not initialized")
	}

	var resp []models.Movie

	if err := m.Movies.FindAll(context.Background(), &resp); err != nil {
		return nil, errors.Wrap(err, "m.Movies.FindAll")
	}

	return resp, nil
}

func (m *Manager) Movie(key string) (models.Movie, error) {
	if m.DB.Database == nil {
		return models.Movie{}, errors.New("database was not initialized")
	}

	var movie models.Movie
	if err := m.Movies.FindByID(context.Background(), &movie, key); err != nil {
		return models.Movie{}, errors.Wrap(err, "m.Movies.FindByID")
	}

	return movie, nil
}
