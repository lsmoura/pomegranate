package models

import (
	"encoding/json"
	"fmt"
	"pomegranate/database"
)

const MovieBucketName = "movies"

type NzbStatus string

const (
	StatusAvailable NzbStatus = "available"
	StatusSnatched            = "snatched"
	StatusFailed              = "failed"
	StatusSuccess             = "success"
	StatusUnknown             = "unknown"
	StatusError               = "error"
)

const MovieKind = "movie"

type NzbInfo struct {
	GUID   string    `json:"guid"`
	ID     string    `json:"id"`
	Size   int64     `json:"size"`
	Status NzbStatus `json:"status"`
	Title  string    `json:"title"`
	URL    string    `json:"url"`

	DownloaderId string `json:"downloader_id"`
}

type Movie struct {
	ImdbId      string    `json:"imdb_id"`
	Title       string    `json:"title"`
	Overview    string    `json:"overview"`
	ReleaseDate string    `json:"release_date"`
	NzbInfo     []NzbInfo `json:"nzb_info"`
}

func (m *Movie) Kind() string {
	return MovieKind
}

func (m *Movie) SetKey(key database.Key) {
	m.ImdbId = string(key)
}

func (m *Movie) GetKey() database.Key {
	return []byte(m.ImdbId)
}

// Store saves the current movie data to the database
func (m Movie) Store(db *database.DB) error {
	dbBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err := db.Store(MovieBucketName, m.GetKey(), dbBytes); err != nil {
		return fmt.Errorf("DB.Store: %w", err)
	}

	return nil
}
