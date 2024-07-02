// Package repo provides interfaces and implementations for managing URLs and events in the repository.
//
// This package includes interfaces for defining repository operations and implementations
// for interacting with databases and in-memory storage.
//
// The main interface provided is Repo, which defines methods for managing URLs and events.
//
// Implementations of Repo interface should be provided for specific database systems or in-memory storage,
// allowing flexibility in choosing the underlying storage mechanism.
package repo

import (
	_ "github.com/jackc/pgx/v5"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
)

// Repo represents the repository interface for managing URLs and events.
type Repo interface {
	// Initialize initializes the repository.
	Initialize() error

	// TableExist checks if the necessary tables exist in the database.
	TableExist() bool

	// Add adds a new URL entry to the repository.
	Add(originalURL, shortURL, userID string) error

	// CreateShortURLs creates multiple short URLs in the repository for the given user.
	CreateShortURLs(urls []*entities.ReqURL, userID string) error

	// Update updates the short URL for the given original URL.
	Update(originalURL, shortURL string) error

	// GetByShortURL retrieves the original URL associated with the given short URL.
	// It returns the original URL and a boolean indicating whether the URL exists.
	GetByShortURL(shortURL string) (string, bool)

	// GetByOriginalURL retrieves the short URL associated with the given original URL.
	// It returns the short URL and a boolean indicating whether the URL exists.
	GetByOriginalURL(originalURL string) (string, bool)

	// ListAll retrieves all short URLs associated with the given user.
	ListAll(userID string) ([]*entities.ReqListAll, error)

	// DeleteURLs deletes multiple URLs associated with the given user.
	DeleteURLs(urls []string, userID string) error

	// URLDeleted checks if the URL with the given short URL is deleted.
	URLDeleted(shortURL string) bool

	// WriteEvent writes an event to the storage.
	WriteEvent(event *inmemory.Event) error

	// SyncData synchronizes data.
	SyncData()

	// CountUsers counts unique users.
	CountUsers() (int, error)

	// CountURLs counts unique URLs.
	CountURLs() (int, error)

	// Ping pings the database to check its availability.
	Ping() error
}
