package themoviedb

type Entry struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIds         []int   `json:"genre_ids"`
	Id               int32   `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type SingleMovieResponse struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection struct {
		Id           int32  `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int32 `json:"budget"`
	Genres []struct {
		Id   int32  `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	Id                  int32   `json:"id"`
	ImdbId              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ProductionCompanies []struct {
		Id            int32  `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int32  `json:"revenue"`
	Runtime         int32  `json:"runtime"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status            string  `json:"status"`
	Tagline           string  `json:"tagline"`
	Title             string  `json:"title"`
	Video             bool    `json:"video"`
	VoteAverage       float64 `json:"vote_average"`
	VoteCount         int32   `json:"vote_count"`
	AlternativeTitles struct {
		Titles []struct {
			Iso31661 string `json:"iso_3166_1"`
			Title    string `json:"title"`
			Type     string `json:"type"`
		} `json:"titles"`
	} `json:"alternative_titles"`
}

type Response struct {
	Page         int32 `json:"page"`
	TotalPages   int32 `json:"total_pages"`
	TotalResults int32 `json:"total_results"`

	Results []Entry
}
