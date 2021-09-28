package themoviedb

type Themoviedb struct {
	apiKey string
}

func New(apiKey string) Themoviedb {
	return Themoviedb{
		apiKey: apiKey,
	}
}

func (d Themoviedb) ReadSingleMovie(key string) (SingleMovieResponse, error) {
	return ReadSingleMovie(d.apiKey, key)
}

func (d Themoviedb) ReadMovies(search string, page int) (Response, error) {
	return ReadMovies(d.apiKey, search, page)
}
