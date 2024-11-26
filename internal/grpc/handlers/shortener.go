package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"github.com/northmule/shorturl/internal/grpc/handlers/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenerHandler хэндлер сокращения ссылок.
type ShortenerHandler struct {
	contract.UnimplementedShortenerHandlerServer
	service *url.ShortURLService
	finder  handlers.Finder
	setter  handlers.LikeURLToUserSetter
}

// NewShortenerHandler конструктор.
func NewShortenerHandler(urlService *url.ShortURLService, finder handlers.Finder, setter handlers.LikeURLToUserSetter) *ShortenerHandler {
	shortenerHandler := &ShortenerHandler{
		service: urlService,
		finder:  finder,
		setter:  setter,
	}
	return shortenerHandler
}

// Shortener обработчик создания короткой ссылки.
func (s *ShortenerHandler) Shortener(ctx context.Context, request *contract.ShortenerRequest) (*contract.ShortenerResponse, error) {

	if !strings.Contains(request.GetUrl(), "http://") && !strings.Contains(request.GetUrl(), "https://") {
		return nil, status.Error(codes.InvalidArgument, "expected url")
	}

	userUUID, err := utils.FillUserUUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "expected userUUID")
	}

	shortURL, err := s.fillShortURL(userUUID, request.GetUrl())
	if err != nil {
		return nil, err
	}

	response := &contract.ShortenerResponse{}
	response.ShortUrl = fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, shortURL)

	return response, nil
}

// ShortenerJSON аналог метода http по сигнатуре ответа
func (s *ShortenerHandler) ShortenerJSON(ctx context.Context, request *contract.ShortenerJSONRequest) (*contract.ShortenerJSONResponse, error) {

	if !strings.Contains(request.GetUrl(), "http://") && !strings.Contains(request.GetUrl(), "https://") {
		return nil, status.Error(codes.InvalidArgument, "\"expected url")
	}

	userUUID, err := utils.FillUserUUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "\"expected userUUID")
	}

	shortURL, err := s.fillShortURL(userUUID, request.GetUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &contract.ShortenerJSONResponse{}
	response.Result = fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, shortURL)

	return response, nil
}

// ShortenerBatch обработка списка адресов.
func (s *ShortenerHandler) ShortenerBatch(ctx context.Context, request *contract.ShortenerBatchRequest) (*contract.ShortenerBatchResponse, error) {

	if len(request.Items) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "expected Items")
	}

	urls := make([]string, 0)
	for _, requestItem := range request.Items {
		if !strings.Contains(requestItem.GetOriginalUrl(), "http://") && !strings.Contains(requestItem.GetOriginalUrl(), "https://") {
			continue
		}
		urls = append(urls, requestItem.GetOriginalUrl())
	}
	if len(urls) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "expected urls")
	}
	modelURLs, err := s.service.DecodeURLs(urls)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseItems := make([]*contract.ShortenerBatchResponse_Item, 0)
	for _, requestItem := range request.Items {
		for _, modelURL := range modelURLs {
			if requestItem.GetOriginalUrl() == modelURL.URL {
				responseItems = append(responseItems, &contract.ShortenerBatchResponse_Item{
					CorrelationId: requestItem.CorrelationId,
					ShortUrl:      fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, modelURL.ShortURL),
				})
			}
		}
	}

	response := &contract.ShortenerBatchResponse{}
	response.Items = responseItems

	return response, nil
}

func (s *ShortenerHandler) fillShortURL(userUUID string, url string) (string, error) {
	var (
		shortURL    string
		isURLExists bool
	)
	shortURLData, err := s.service.DecodeURL(url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.CodeErrorDuplicateKey {
			isURLExists = true
		} else {
			return "", status.Error(codes.Internal, err.Error())
		}
	}
	if isURLExists {
		modelURL, err := s.finder.FindByURL(url)
		if err != nil {
			return "", status.Error(codes.Internal, err.Error())
		}

		return modelURL.ShortURL, status.Errorf(codes.AlreadyExists, "url already exists")
	}
	shortURL = shortURLData.ShortURL
	err = s.setter.LikeURLToUser(shortURLData.URLID, userUUID)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return shortURL, nil
}
