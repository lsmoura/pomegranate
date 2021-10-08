package sabnzbd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type QueueSlot struct {
	Status     string   `json:"status"`
	Index      int      `json:"index"`
	Password   string   `json:"password"`
	AvgAge     string   `json:"avg_age"`
	Script     string   `json:"script"`
	HasRating  bool     `json:"has_rating"`
	Mb         string   `json:"mb"`
	Mbleft     string   `json:"mbleft"`
	Mbmissing  string   `json:"mbmissing"`
	Size       string   `json:"size"`
	Sizeleft   string   `json:"sizeleft"`
	Filename   string   `json:"filename"`
	Labels     []string `json:"labels"` // Labels like DUPLICATE, ENCRYPTED and PROPAGATING X min
	Priority   string   `json:"priority"`
	Cat        string   `json:"cat"`
	Eta        string   `json:"eta"`
	Timeleft   string   `json:"timeleft"`
	Percentage string   `json:"percentage"`
	NzoId      string   `json:"nzo_id"`
	Unpackopts string   `json:"unpackopts"`
}

type Queue struct {
	Status          string      `json:"status"`
	Speedlimit      string      `json:"speedlimit"`     // In percentage of maximum set by user
	SpeedlimitAbs   string      `json:"speedlimit_abs"` // In bytes/s
	Paused          bool        `json:"paused"`
	NoofslotsTotal  int         `json:"noofslots_total"`
	Noofslots       int         `json:"noofslots"`
	Limit           int         `json:"limit"`
	Start           int         `json:"start"`
	Eta             string      `json:"eta"`
	Timeleft        string      `json:"timeleft"`
	Speed           string      `json:"speed"`
	Kbpersec        string      `json:"kbpersec"`
	Size            string      `json:"size"`
	Sizeleft        string      `json:"sizeleft"`
	Mb              string      `json:"mb"`
	Mbleft          string      `json:"mbleft"`
	Slots           []QueueSlot `json:"slots"`
	Categories      []string    `json:"categories"`
	Scripts         []string    `json:"scripts"`
	Diskspace1      string      `json:"diskspace1"`
	Diskspace2      string      `json:"diskspace2"`
	Diskspacetotal1 string      `json:"diskspacetotal1"`
	Diskspacetotal2 string      `json:"diskspacetotal2"`
	Diskspace1Norm  string      `json:"diskspace1_norm"`
	Diskspace2Norm  string      `json:"diskspace2_norm"`
	RatingEnable    bool        `json:"rating_enable"`
	HaveWarnings    string      `json:"have_warnings"`
	PauseInt        string      `json:"pause_int"`
	Loadavg         string      `json:"loadavg"` // On Linux this will contain a string with information about system load
	LeftQuota       string      `json:"left_quota"`
	RefreshRate     string      `json:"refresh_rate"`
	Version         string      `json:"version"`
	Finish          int         `json:"finish"`
	CacheArt        string      `json:"cache_art"`
	CacheSize       string      `json:"cache_size"`
	CacheMax        string      `json:"cache_max"`
	Finishaction    interface{} `json:"finishaction"`
	PausedAll       bool        `json:"paused_all"`
	Quota           string      `json:"quota"`
	HaveQuota       bool        `json:"have_quota"`
	QueueDetails    string      `json:"queue_details"`
}

type QueueRequestParams struct {
	Start  int32   `json:"start"`   // Index of job to start at
	Limit  int32   `json:"limit"`   // Number of jobs to display
	Search string  `json:"search"`  // Filter job names by search term
	NzoIds []int32 `json:"nzo_ids"` // Filter jobs by nzo_ids
}

type QueueResponse struct {
	Queue Queue `json:"queue"`
}

type PauseQueueResponse struct {
	Status bool `json:"status"`
}

func (s Sabnzbd) Queue(params QueueRequestParams) (Queue, error) {
	if s.Host == "" {
		return Queue{}, errors.New("sabnzbd structure has no host")
	}

	// TODO: detect if scheme is present on Host
	u := s.url()
	query := u.Query()
	query.Add("mode", "queue")
	if params.Limit != 0 {
		query.Add("limit", strconv.Itoa(int(params.Limit)))
	}
	if len(params.NzoIds) > 0 {
		var ids []string
		for _, id := range params.NzoIds {
			ids = append(ids, strconv.Itoa(int(id)))
		}
		query.Add("nzo_ids", strings.Join(ids, ","))
	}
	if params.Start > 0 {
		query.Add("start", strconv.Itoa(int(params.Start)))
	}
	if params.Search != "" {
		query.Add("search", params.Search)
	}
	u.RawQuery = query.Encode()

	s.log("HTTP get: %s\n", strings.ReplaceAll(u.String(), s.Apikey, "xxx"))
	resp, err := http.Get(u.String())
	if err != nil {
		return Queue{}, errors.Wrap(err, "http.Get")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("Body.Close: %w", err))
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Queue{}, errors.Wrap(err, "ioutil.ReadAll")
	}

	var queue QueueResponse
	if err := json.Unmarshal(body, &queue); err != nil {
		return Queue{}, errors.Wrap(err, "json.Unmarshal")
	}

	return queue.Queue, nil
}

func (s Sabnzbd) PauseQueue() error {
	u := s.url()

	query := u.Query()
	query.Add("mode", "pause")
	u.RawQuery = query.Encode()

	s.log("HTTP get: %s\n", strings.ReplaceAll(u.String(), s.Apikey, "xxx"))
	resp, err := http.Get(u.String())
	if err != nil {
		return errors.Wrap(err, "http.Get")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("Body.Close: %w", err))
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll")
	}

	var apiResponse PauseQueueResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	if !apiResponse.Status {
		return errors.New("response status is false")
	}

	return nil
}

func (s Sabnzbd) ResumeQueue() error {
	u := s.url()

	query := u.Query()
	query.Add("mode", "resume")
	u.RawQuery = query.Encode()

	fmt.Printf("HTTP get (sabnzbd): %s\n", strings.ReplaceAll(u.String(), s.Apikey, "xxx"))
	resp, err := http.Get(u.String())
	if err != nil {
		return errors.Wrap(err, "http.Get")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("Body.Close: %w", err))
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll")
	}

	var apiResponse PauseQueueResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	if !apiResponse.Status {
		return errors.New("response status is false")
	}

	return nil
}
