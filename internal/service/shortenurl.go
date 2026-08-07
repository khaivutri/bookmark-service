package service

import (
	"context"
	"errors"
	"time"

	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/khaivutri/bookmark-service/pkg/utils"
	
)

// ShortenURL defines the interface for URL shortening operations.
type ShortenURL interface {
	// CreateCodeFromLink generates a unique short code for a URL and stores it in Redis.
	CreateCodeFromLink(ctx context.Context, url string, exp int64) (string, error)
	// GetLinkFromCode retrieves the original URL corresponding to the short code.
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

type shortenURL struct {
	storage 	repository.URLStorage
	generator 	utils.GenCode
}	

// NewURLStorage constructs a new shortenURL service instance.
func NewURLStorage(storage repository.URLStorage, generator utils.GenCode) ShortenURL {
	return &shortenURL{storage: storage, generator: generator}
}

const CODE_LEN =7
// CreateCodeFromLink generates a unique short code for a URL and stores it in Redis.
func (s *shortenURL) CreateCodeFromLink(ctx context.Context, url string, exp int64) (string, error){

	code, errGen := s.generator.Generate(CODE_LEN)
	if errGen != nil {
		return "", errGen
	}

	result, errGet := s.storage.GetURL(ctx, code)
	if errGet != nil && !errors.Is(errGet, repository.ErrCodeNotFound) {
		return "", errGet
	}

	if result != "" {
		return s.CreateCodeFromLink(ctx, url, exp)
	}

	errSto := s.storage.StoreURL(ctx, code, url, time.Duration(exp)*time.Second)
	if errSto != nil {
		return "", errSto
	}
	
	return code, nil
}

// GetLinkFromCode retrieves the original URL corresponding to the short code.
func (s *shortenURL) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return s.storage.GetURL(ctx, code)
}