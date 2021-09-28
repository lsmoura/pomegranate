package database

import (
	"encoding/json"
	"fmt"
)

type NzbStatus string

const (
	StatusAvailable NzbStatus = "available"
	StatusSnatched            = "snatched"
	StatusFailed              = "failed"
	StatusSuccess             = "success"
	StatusUnknown             = "unknown"
	StatusError               = "error"
)

type NzbInfo struct {
	Title  string    `json:"title"`
	Id     string    `json:"id"`
	URL    string    `json:"url"`
	Status NzbStatus `json:"status"`
}

type Movie struct {
	ImdbId      string    `json:"imdb_id"`
	Title       string    `json:"title"`
	Overview    string    `json:"overview"`
	ReleaseDate string    `json:"release_date"`
	NzbInfo     []NzbInfo `json:"nzb_info"`
}

func (m Movie) Store(db DB) error {
	dbBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err := db.Store(MovieBucketName, []byte(m.ImdbId), dbBytes); err != nil {
		return fmt.Errorf("DB.Store: %w", err)
	}

	return nil
}
