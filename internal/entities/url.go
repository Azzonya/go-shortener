// Package entities provides data structures used across the URL shortener application.
package entities

// ReqURL represents a request structure for URL shortening.
type ReqURL struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url,omitempty"`
	ShortURL    string `json:"short_url"`
}

// ReqListAll represents a structure for listing all shortened URLs.
type ReqListAll struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Storage represents a structure for storing URLs in a database.
type Storage struct {
	UUID        string `db:"id"`
	ShortURL    string `db:"shorturl"`
	OriginalURL string `db:"originalurl"`
	UserID      string `db:"userid"`
	DeletedFlag bool   `db:"deleted"`
}
