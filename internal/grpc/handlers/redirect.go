package handlers

import (
	"context"

	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RedirectHandler хэндлер для обработки коротких ссылок.
type RedirectHandler struct {
	contract.UnimplementedRedirectHandlerServer
	service *url.ShortURLService
}

// NewRedirectHandler конструктор хэндлера.
func NewRedirectHandler(urlService *url.ShortURLService) *RedirectHandler {
	redirectHandler := &RedirectHandler{
		service: urlService,
	}
	return redirectHandler
}

// Redirect обработчик получения оригинальной ссылки из короткой.
func (r *RedirectHandler) Redirect(ctx context.Context, request *contract.RedirectRequest) (*contract.RedirectResponse, error) {

	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "expected id value")
	}
	modelURL, err := r.service.EncodeShortURL(request.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}
	if !modelURL.DeletedAt.IsZero() {
		return nil, status.Error(codes.NotFound, "expected id value")
	}

	response := &contract.RedirectResponse{}
	response.Url = modelURL.URL

	return response, nil
}
