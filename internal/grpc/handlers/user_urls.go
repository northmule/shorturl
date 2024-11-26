package handlers

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"github.com/northmule/shorturl/internal/grpc/handlers/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserURLsHandler хэндлер отображения ссылок пользователя.
type UserURLsHandler struct {
	contract.UnimplementedUserUrlsHandlerServer
	finder  handlers.URLFinder
	session storage.SessionAdapter
	worker  handlers.Deleter
}

// NewUserURLsHandler Конструктор.
func NewUserURLsHandler(finder handlers.URLFinder, sessionStorage storage.SessionAdapter, worker handlers.Deleter) *UserURLsHandler {
	instance := UserURLsHandler{
		finder:  finder,
		session: sessionStorage,
		worker:  worker,
	}
	return &instance
}

// View короткие ссылки пользователя.
func (u *UserURLsHandler) View(ctx context.Context, request *contract.ViewRequest) (*contract.ViewResponse, error) {
	userUUID, err := utils.FillUserUUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "expected userUUID")
	}
	userURLs, err := u.finder.FindUrlsByUserID(userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(*userURLs) == 0 {
		return nil, status.Error(codes.NotFound, "url not found")
	}

	var responseList []*contract.ViewResponse_Item
	for _, urlItem := range *userURLs {
		responseList = append(responseList, &contract.ViewResponse_Item{
			ShortUrl:    fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, urlItem.ShortURL),
			OriginalUrl: urlItem.URL,
		})
	}

	response := &contract.ViewResponse{}
	response.Items = responseList

	return response, nil
}

// Delete удаление ссылок текущего пользователя.
func (u *UserURLsHandler) Delete(ctx context.Context, request *contract.DeleteRequest) (*empty.Empty, error) {

	userUUID, err := utils.FillUserUUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "expected userUUID")
	}

	u.worker.Del(userUUID, request.GetShortUrls())

	response := &empty.Empty{}

	return response, nil
}
