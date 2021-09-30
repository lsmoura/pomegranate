package service

import (
	"fmt"
	"log"
	"net/http"
	"pomegranate/database"
	"pomegranate/manager"
	"pomegranate/newznab"
	"pomegranate/sabnzbd"
	"pomegranate/themoviedb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	DB      database.DB
	Newz    []newznab.Newznab
	Sabnzbd sabnzbd.Sabnzbd
	Tmdb    themoviedb.Themoviedb
}

type MovieSearchResponse struct {
	Movies []manager.MovieEntry `json:"movies"`
}

func internalError(w http.ResponseWriter, format string, a ...interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("internal error"))
	if err != nil {
		log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
	}
	log.Println(fmt.Errorf(format, a...))
}

func Service(config Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pomegranate"))
		if err != nil {
			log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
		}
	})
	r.Get("/movie/search", config.movieSearchHandler)
	r.Get("/movie/add", config.movieAddHandler)
	r.Get("/movie/list", config.movieListHandler)

	r.Get("/nzb/download", config.nzbDownload)

	return r
}
