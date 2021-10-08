package sabnzbd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Logger interface {
	Log(serviceName string, format string, a ...interface{})
}

type Sabnzbd struct {
	Host   string
	Apikey string
	Logger Logger
}

func New(host string, apikey string) Sabnzbd {
	return Sabnzbd{
		Host:   host,
		Apikey: apikey,
	}
}

func (s Sabnzbd) log(format string, a ...interface{}) {
	if s.Logger == nil {
		return
	}

	s.Logger.Log("sabnzbd", format, a...)
}

func (s Sabnzbd) url() *url.URL {
	// TODO: detect if scheme is present on Host
	u := new(url.URL)

	u.Scheme = "http"
	u.Host = s.Host
	u.Path = "api"

	query := u.Query()
	if s.Apikey != "" {
		query.Add("apikey", s.Apikey)
	}
	query.Add("output", "json")
	u.RawQuery = query.Encode()

	return u
}

type PriorityType int32

const (
	PriorityDefault   PriorityType = -100
	PriorityDuplicate              = -3
	PriorityPaused                 = -2
	PriorityLow                    = -1
	PriorityNormal                 = 0
	PriorityHigh                   = 1
	PriorityForce                  = 2
)

type AddUrlParams struct {
	Name     string       `query_name:"name"`     // link to the NZB to be fetched.
	NzbName  string       `query_name:"nzbname"`  // name of the job, if empty the NZB filename is used.
	Password string       `query_name:"password"` // password to use when unpacking the job.
	Category string       `query_name:"cat"`      // category to be assigned, * means Default. List of available categories can be retrieved from get_cats.
	Script   string       `query_name:"script"`   // script to be assigned, Default will use the script assigned to the category. List of available scripts can be retrieved from get_scripts.
	Priority PriorityType `query_name:"priority"` // priority to be assigned
	PP       int32        `query_name:"pp"`       // post-processing options
}

type AddUrlResponse struct {
	Status bool     `json:"status"`
	NzoIds []string `json:"nzo_ids"`
}

func (s Sabnzbd) AddUrl(params AddUrlParams) ([]string, error) {
	if params.Name == "" {
		return nil, fmt.Errorf("name is required for adding a url")
	}
	u := s.url()
	query := u.Query()
	query.Set("mode", "addurl")
	if err := InjectQuery(query, params); err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	fmt.Printf("HTTP get (sabnzbd): %s\n", strings.ReplaceAll(u.String(), s.Apikey, "xxx"))
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "http.Get")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("Body.Close: %w", err))
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	var apiResponse AddUrlResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	if !apiResponse.Status {
		fmt.Println(string(body))
		return nil, errors.New("response status is false")
	}

	return apiResponse.NzoIds, nil
}
