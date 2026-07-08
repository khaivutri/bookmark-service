package service

import (
	"context"
	"errors"
	"time"

	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/khaivutri/bookmark-service/pkg/utils"
	
)

type ShortenURL interface {
	CreateCodeFromLink(ctx context.Context, url string, exp int64) (string, error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

type shortenURL struct {
	storage 	repository.URLStorage
	generator 	utils.GenCode
}	

func NewURLStorage(storage repository.URLStorage, generator utils.GenCode) ShortenURL {
	return &shortenURL{storage: storage, generator: generator}
}

const CODE_LEN =7

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

	errSto := s.storage.StoreURL(ctx, code, url, time.Duration(exp))
	if errSto != nil {
		return "", errSto
	}
	
	return code, nil
}


func (s *shortenURL) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return s.storage.GetURL(ctx, code)
}