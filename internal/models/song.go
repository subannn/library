package models

type Song struct {
	Group    string `json:"group"`
	SongName string `json:"song"`
}

type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongQuery struct {
	Group       string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongResponse struct {
	SongID      int64  `db:"song_id" json:"song_id"`
	GroupID     int64  `db:"group_id" json:"group_id"`
	SongName    string `db:"song_name" json:"song_name"`
	GroupName   string `db:"group_name" json:"group_name"`
	ReleaseDate string `db:"release_date" json:"release_date"`
}
