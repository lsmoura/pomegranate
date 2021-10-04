package service

import (
	"encoding/json"
	"fmt"
	"github.com/lsmoura/humantoken"
	"log"
	"net/http"
	"pomegranate/database"
)

type MovieAddResponse struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Overview string `json:"overview"`
}

func (c Config) movieSearchHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")

	movies, err := c.Manager.MovieSearch(c.Tmdb, searchQuery)
	if err != nil {
		internalError(w, "manager.MovieSearch: %w", err)
		return
	}

	payload := MovieSearchResponse{
		Movies: movies,
	}

	w.Header().Add("Content-Type", "application/json")
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		internalError(w, "json.Marshal: %w", err)
		return
	}

	if _, err := w.Write(payloadBytes); err != nil {
		log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
	}
}

func (c Config) movieAddHandler(w http.ResponseWriter, r *http.Request) {
	identifier := r.URL.Query().Get("identifier")

	// TODO: Validate that the identifier is in the format tt0000000...
	movie, err := c.Tmdb.ReadSingleMovie(identifier)
	if err != nil {
		internalError(w, "tmdb.ReadSingleMovie (%s): %w", movie.Id, err)
		return
	}

	dbMovie, err := c.DB.Movie(identifier)
	if err != nil {
		internalError(w, "DB.Movie (%s): %w", identifier, err)
		return
	}

	dbMovie.ImdbId = movie.ImdbId
	dbMovie.Title = movie.Title
	dbMovie.ReleaseDate = movie.ReleaseDate
	dbMovie.Overview = movie.Overview

	for _, n := range c.Newz {
		parsedIdentifier := identifier[2:]

		items, err := n.SearchImdb(parsedIdentifier)
		if err != nil {
			internalError(w, "newznab.SearchImdb: %w", err)
			return
		}

		fmt.Println(items)

		for _, item := range items {
			found := false
			for _, nzb := range dbMovie.NzbInfo {
				if nzb.URL == item.URL {
					found = true
				}
			}
			if !found {
				dbMovie.NzbInfo = append(dbMovie.NzbInfo, database.NzbInfo{
					ID:     humantoken.Generate(8, nil),
					Title:  item.Title,
					GUID:   item.GUID,
					URL:    item.URL,
					Status: database.StatusUnknown,
					Size:   item.Size,
				})
			}
		}
	}

	if err := dbMovie.Store(c.DB); err != nil {
		internalError(w, "database.Movie.Store: %w", err)
		return
	}

	response := MovieAddResponse{
		Message:  "Movie added",
		Title:    dbMovie.Title,
		Overview: dbMovie.Overview,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		internalError(w, "json.Marshal: %w", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte(responseBytes)); err != nil {
		log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
	}
}

func (c Config) movieListHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := c.DB.AllMovies()
	if err != nil {
		internalError(w, "DB.AllMovies: %w", err)
		return
	}

	moviesBytes, err := json.Marshal(movies)
	if err != nil {
		internalError(w, "json.Marshal: %w", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte(moviesBytes)); err != nil {
		log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
	}
}
