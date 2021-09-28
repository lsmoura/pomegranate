package service

import (
	"encoding/json"
	"fmt"
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

	// TODO: Error if query is empty

	res, err := c.Tmdb.ReadMovies(searchQuery, 0)
	if err != nil {
		internalError(w, "tmdb.ReadMovies: %w", err)
		return
	}

	payload := MovieSearchResponse{}

	for _, movie := range res.Results {
		m := MovieEntry{
			Titles:   []string{movie.Title},
			Released: movie.ReleaseDate,
			TmdbId:   movie.Id,
		}

		extraInfo, err := c.Tmdb.ReadSingleMovie(fmt.Sprintf("%d", movie.Id))
		if err != nil {
			internalError(w, "tmdb.ReadSingleMovie (%d): %w", movie.Id, err)
			return
		}

		m.ImdbId = extraInfo.ImdbId
		m.Runtime = extraInfo.Runtime

		payload.Movies = append(payload.Movies, m)
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
	// TODO: Check if the movie already exists and update

	movie, err := c.Tmdb.ReadSingleMovie(identifier)
	if err != nil {
		internalError(w, "tmdb.ReadSingleMovie (%s): %w", movie.Id, err)
		return
	}

	dbMovie := database.Movie{
		ImdbId:      movie.ImdbId,
		Title:       movie.Title,
		ReleaseDate: movie.ReleaseDate,
		Overview:    movie.Overview,
	}

	if err := dbMovie.Store(c.DB); err != nil {
		internalError(w, "database.Movie.Store: %w", err)
		return
	}

	parsedIdentifier := identifier[2:]

	for _, n := range c.Newz {
		items, err := n.SearchImdb(parsedIdentifier)
		if err != nil {
			internalError(w, "newznab.SearchImdb: %w", err)
			return
		}

		fmt.Println(items)

		// TODO: Update database
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
