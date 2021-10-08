package themoviedb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// ReadSingleMovie takes a movie key as parameter.
// the key parameter can be either a themoviedb id or an imdb id (which beging with tt)
func ReadSingleMovie(apiKey string, key string, logger func(string, ...interface{})) (SingleMovieResponse, error) {
	movieInfo := SingleMovieResponse{}
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s&append_to_response=alternative_titles&language=en", key, apiKey)

	logger(fmt.Sprintf("HTTP request: %s\n", strings.ReplaceAll(url, apiKey, "xxx")))
	resp, err := http.Get(url)
	if err != nil {
		return movieInfo, fmt.Errorf("http.Get: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger(fmt.Sprintf("Body.Close: %w", err))
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return movieInfo, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	if err := json.Unmarshal(body, &movieInfo); err != nil {
		return movieInfo, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return movieInfo, nil
}

func ReadMovies(apiKey string, search string, page int, logger func(string, ...interface{})) (Response, error) {
	ThemoviedbResponse := Response{}

	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&search_type=ngram&query=%s", apiKey, search)
	if page != 0 {
		url = fmt.Sprintf("%s&page=%d", url, page)
	}

	logger(fmt.Sprintf("HTTP request: %s\n", strings.ReplaceAll(url, apiKey, "xxx")))
	resp, err := http.Get(url)
	if err != nil {
		return ThemoviedbResponse, fmt.Errorf("http.Get: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger(fmt.Sprintf("Body.Close: %w", err))
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ThemoviedbResponse, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	if err := json.Unmarshal(body, &ThemoviedbResponse); err != nil {
		return ThemoviedbResponse, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return ThemoviedbResponse, nil
}
