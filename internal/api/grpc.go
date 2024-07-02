package api

import (
	"context"
	"errors"
	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/session"
	shortener_service "github.com/Azzonya/go-shortener/internal/shortener"
	pb "github.com/Azzonya/go-shortener/pkg/proto/shortener"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

type St struct {
	pb.UnimplementedShortenerServer
	shortener *shortener_service.Shortener
}

func NewGrpcHandlers(shortener *shortener_service.Shortener) *St {
	return &St{
		shortener: shortener,
	}
}

// Shorten processes a request to shorten a URL and returns the shortened URL or an error.
//
// Parameters:
//   - ctx: The context of the request.
//   - req: The ShortenRequest containing the URL to shorten.
//
// Returns:
//   - *pb.ShortenResponse: The response containing the shortened URL.
//   - error: An error if the URL couldn't be shortened.
func (s *St) Shorten(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	var userID string
	var exist bool

	reqObj := strings.TrimSpace(req.GetUrl())
	if reqObj == "" {
		return nil, errors.New("url is empty")
	}

	user, ok := session.GetUserFromMetadata(ctx)
	if ok {
		userID = user.ID
	}

	outputURL, err := s.shortener.ShortenAndSaveLink(reqObj, userID)
	if err != nil {
		outputURL, exist = s.shortener.GetOneByOriginalURL(reqObj)
		if !exist {
			return &pb.ShortenResponse{ShortenedUrl: outputURL}, errors.New("failed to add URL to in-memory storage")
		}
	}

	return &pb.ShortenResponse{ShortenedUrl: outputURL}, nil
}

// Redirect redirects a request to the original URL corresponding to a given short URL.
//
// Parameters:
//   - ctx: The context of the request.
//   - req: The RedirectRequest containing the short URL to redirect.
//
// Returns:
//   - *pb.RedirectResponse: The response containing the original URL.
//   - error: An error if the original URL couldn't be retrieved.
func (s *St) Redirect(_ context.Context, req *pb.RedirectRequest) (*pb.RedirectResponse, error) {
	shortURL := req.GetShortUrl()
	if shortURL == "" {
		return nil, errors.New("short_url is empty")
	}

	if s.shortener.IsDeleted(shortURL) {
		return nil, errors.New("short_url is deleted")
	}

	originalURL, exist := s.shortener.GetOneByShortURL(shortURL)
	if !exist {
		return nil, errors.New("failed to get original URL")
	}

	return &pb.RedirectResponse{OriginalUrl: originalURL}, nil
}

// ShortenURLs processes a request to shorten multiple URLs and returns the shortened URLs or an error.
//
// Parameters:
//   - ctx: The context of the request.
//   - req: The ShortenURLsRequest containing the URLs to shorten.
//
// Returns:
//   - *pb.ShortenURLsResponse: The response containing the shortened URLs.
//   - error: An error if the URLs couldn't be shortened.
func (s *St) ShortenURLs(ctx context.Context, req *pb.ShortenURLsRequest) (*pb.ShortenURLsResponse, error) {
	var userID string

	user, ok := session.GetUserFromContext(ctx)
	if ok {
		userID = user.ID
	}

	var urls []*entities.ReqURL
	for _, u := range req.GetUrls() {
		urls = append(urls, &entities.ReqURL{OriginalURL: u.GetUrl()})
	}

	shortenedURLs, err := s.shortener.ShortenURLs(urls, userID)
	if err != nil {
		return nil, err
	}

	var shortenedURLStrings []string
	for _, su := range shortenedURLs {
		shortenedURLStrings = append(shortenedURLStrings, su.ShortURL)
	}

	return &pb.ShortenURLsResponse{ShortenedUrls: shortenedURLStrings}, nil
}

// ListAll retrieves all URLs associated with the authenticated user.
//
// Parameters:
//   - ctx: The context of the request.
//   - _: The Empty message (unused).
//
// Returns:
//   - *pb.ListAllRequest: The response containing the list of URLs.
//   - error: An error if the URLs couldn't be retrieved.
func (s *St) ListAll(ctx context.Context, _ *emptypb.Empty) (*pb.ListAllRequest, error) {
	user, ok := session.GetUserFromContext(ctx)
	if !ok {
		return nil, errors.New("failed to get user")
	}

	if user.IsNew() {
		return nil, errors.New("user is new")
	}

	result, err := s.shortener.ListAll(user.ID)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &pb.ListAllRequest{}, nil
	}

	var listURLs []*pb.List
	for _, su := range result {
		listURLs = append(listURLs, &pb.List{
			OriginalUrl: su.OriginalURL,
			ShortUrl:    su.ShortURL,
		})
	}

	return &pb.ListAllRequest{URLs: listURLs}, nil
}

// Ping checks the availability of the service.
//
// Parameters:
//   - _: The Empty message (unused).
//
// Returns:
//   - *emptypb.Empty: An empty response indicating the service is available.
//   - error: An error if the ping request failed.
func (s *St) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// DeleteURLs deletes multiple URLs associated with the authenticated user.
//
// Parameters:
//   - ctx: The context of the request.
//   - req: The DeleteURLsRequest containing the short URLs to delete.
//
// Returns:
//   - *emptypb.Empty: An empty response indicating the URLs are deleted.
//   - error: An error if the URLs couldn't be deleted.
func (s *St) DeleteURLs(ctx context.Context, req *pb.DeleteURLsRequest) (*emptypb.Empty, error) {
	user, ok := session.GetUserFromContext(ctx)
	if !ok {
		return nil, errors.New("failed to get user")
	}

	var shortURLs []string
	for i := 0; i < len(req.ShortUrls); i++ {
		shortURLs = append(shortURLs, req.ShortUrls[i].ShortUrl)
	}

	s.shortener.DeleteURLs(shortURLs, user.ID)

	return &emptypb.Empty{}, nil
}

// Stats retrieves statistics about the URLs and users managed by the service.
//
// Parameters:
//   - _: The Empty message (unused).
//
// Returns:
//   - *pb.StatsResponse: The response containing the statistics.
//   - error: An error if the statistics couldn't be retrieved.
func (s *St) Stats(_ context.Context, _ *emptypb.Empty) (*pb.StatsResponse, error) {
	stats, err := s.shortener.GetStats()
	if err != nil {
		return nil, errors.New("failed to get stats")
	}

	return &pb.StatsResponse{
		URLs:  int32(stats.URLs),
		Users: int32(stats.Users),
	}, nil
}
