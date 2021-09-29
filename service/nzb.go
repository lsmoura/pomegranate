package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pomegranate/database"
	"pomegranate/sabnzbd"
)

func (c Config) nzbDownload(w http.ResponseWriter, r *http.Request) {
	nzbID := r.URL.Query().Get("id")

	// TODO: Error if id is empty

	movie, err := c.DB.MovieWithNzbID(nzbID)
	if err != nil {
		internalError(w, "database.MovieWithNzbID: %w", err)
		return
	}

	if movie.Title == "" {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("not found")); err != nil {
			log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
		}
		return
	}

	var nzb *database.NzbInfo
	for _, info := range movie.NzbInfo {
		if info.ID == nzbID {
			nzb = &info
		}
	}
	if nzb == nil {
		log.Printf("nzb id %s not found. this error should not be possible\n", nzbID)
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("not found")); err != nil {
			log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
		}
		return
	}

	ids, err := c.Sabnzbd.AddUrl(sabnzbd.AddUrlParams{Name: nzb.URL})
	if err != nil {
		internalError(w, "Sabnzbd.AddUrl: %w", err)
		return
	}
	if len(ids) > 1 {
		log.Printf("I don't know what to do with this many ids! %s\n", ids)
	}
	if len(ids) < 1 {
		internalError(w, "Sabnzbd.AddUrl returned no ids")
		return
	}

	nzb.DownloaderId = ids[0]
	nzb.Status = database.StatusSnatched

	// update Movie nzbs
	var nzbList []database.NzbInfo
	for _, info := range movie.NzbInfo {
		if info.ID != nzb.ID {
			nzbList = append(nzbList, info)
		} else {
			nzbList = append(nzbList, *nzb)
		}
	}
	movie.NzbInfo = nzbList

	if err := movie.Store(c.DB); err != nil {
		internalError(w, "movie.Store: %w", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	payloadBytes, err := json.Marshal(movie)
	if err != nil {
		internalError(w, "json.Marshal: %w", err)
		return
	}

	if _, err := w.Write(payloadBytes); err != nil {
		log.Println(fmt.Errorf("http.ResponseWriter.Write: %w", err))
	}
}
