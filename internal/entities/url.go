package entities

type ReqURL struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url,omitempty"`
	ShortURL    string `json:"short_url"`
}

type ReqURLs struct {
	urls *ReqURL
}
