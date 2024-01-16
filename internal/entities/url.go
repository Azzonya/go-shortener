package entities

type ReqURL struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url,omitempty"`
	ShortURL    string `json:"short_url"`
}

type ReqListAll struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage struct {
	UUID        string `db:"id"`
	ShortURL    string `db:"shorturl"`
	OriginalURL string `db:"originalurl"`
	UserID      string `db:"userid"`
	DeletedFlag bool   `db:"deleted"`
}
