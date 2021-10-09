package themoviedb

type Themoviedb struct {
	apiKey string
	Logger Logger
}

type Logger interface {
	Log(serviceName string, format string, a ...interface{})
}

func (d Themoviedb) log(format string, a ...interface{}) {
	if d.Logger == nil {
		return
	}

	d.Logger.Log("themoviedb", format, a...)
}

func New(apiKey string) Themoviedb {
	return Themoviedb{
		apiKey: apiKey,
	}
}

func (d Themoviedb) ReadSingleMovie(key string) (SingleMovieResponse, error) {
	return ReadSingleMovie(d.apiKey, key, d.log)
}

func (d Themoviedb) ReadMovies(search string, page int) (Response, error) {
	return ReadMovies(d.apiKey, search, page, d.log)
}
