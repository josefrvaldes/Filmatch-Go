package model

type Content struct {
	TMDBID           uint    `json:"id"`
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIDs         []int   `gorm:"-" json:"genre_ids"`
	GenreIDsRaw      string  `json:"-"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}
