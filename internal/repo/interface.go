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
	"context"
	_ "github.com/jackc/pgx/v5"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
)

// Repo represents the repository interface for managing URLs and events.
type Repo interface {
	// Initialize initializes the repository.
	Initialize(ctx context.Context) error

	// TableExist checks if the necessary tables exist in the database.
	TableExist(ctx context.Context) bool

	// Add adds a new URL entry to the repository.
	Add(ctx context.Context, originalURL, shortURL, userID string) error

	// CreateShortURLs creates multiple short URLs in the repository for the given user.
	CreateShortURLs(ctx context.Context, urls []*entities.ReqURL, userID string) error

	// Update updates the short URL for the given original URL.
	Update(ctx context.Context, originalURL, shortURL string) error

	// GetByShortURL retrieves the original URL associated with the given short URL.
	// It returns the original URL and a boolean indicating whether the URL exists.
	GetByShortURL(ctx context.Context, shortURL string) (string, bool)

	// GetByOriginalURL retrieves the short URL associated with the given original URL.
	// It returns the short URL and a boolean indicating whether the URL exists.
	GetByOriginalURL(ctx context.Context, originalURL string) (string, bool)

	// ListAll retrieves all short URLs associated with the given user.
	ListAll(ctx context.Context, userID string) ([]*entities.ReqListAll, error)

	// DeleteURLs deletes multiple URLs associated with the given user.
	DeleteURLs(ctx context.Context, urls []string, userID string) error

	// URLDeleted checks if the URL with the given short URL is deleted.
	URLDeleted(ctx context.Context, shortURL string) bool

	// WriteEvent writes an event to the storage.
	WriteEvent(event *inmemory.Event) error

	// SyncData synchronizes data.
	SyncData()

	// CountUsers counts unique users.
	CountUsers(ctx context.Context) (int, error)

	// CountURLs counts unique URLs.
	CountURLs(ctx context.Context) (int, error)

	// Ping pings the database to check its availability.
	Ping(ctx context.Context) error
}
