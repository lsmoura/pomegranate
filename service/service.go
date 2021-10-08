package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pomegranate/database"
	"pomegranate/manager"
	"pomegranate/newznab"
	"pomegranate/sabnzbd"
	"pomegranate/themoviedb"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	DB      *database.DB
	Newz    []newznab.Newznab
	Sabnzbd sabnzbd.Sabnzbd
	Tmdb    themoviedb.Themoviedb

	Manager *manager.Manager
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

func writeJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if _, err := w.Write(payloadBytes); err != nil {
		return fmt.Errorf("w.Write: %w", err)
	}

	return nil
}

// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	fmt.Printf("filesystem %s\n", path)
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func Service(config Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/movie/search", config.movieSearchHandler)
	r.Get("/movie/add", config.movieAddHandler)
	r.Get("/movie/list", config.movieListHandler)

	r.Get("/nzb/download", config.nzbDownload)

	// TODO: read the filesystem root dir from config
	fileServer := http.FileServer(FileSystem{http.Dir("./static")})
	r.Handle("/*", fileServer)

	return r
}
