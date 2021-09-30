package manager

import (
	"fmt"
	"pomegranate/themoviedb"
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

func MovieSearch(tmdb themoviedb.Themoviedb, query string) ([]MovieEntry, error) {
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
