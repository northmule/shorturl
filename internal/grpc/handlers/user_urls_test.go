package handlers

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/northmule/shorturl/internal/app/workers"
	"github.com/northmule/shorturl/internal/grpc/contract"
	mData "github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MockFinder struct {
	mock.Mock
}

func (m *MockFinder) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	args := m.Called(userUUID)
	return args.Get(0).(*[]models.URL), args.Error(1)
}

type MockFinderBad struct {
	mock.Mock
}

func (m *MockFinderBad) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, errors.New("error")
}

type MockDeleter struct {
	mock.Mock
	IsDeleted bool
}

func (w *MockDeleter) Del(userUUID string, input []string) {
	w.IsDeleted = true
}

func TestUserURLsHandler_View(t *testing.T) {
	_ = logger.InitLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	tests := []struct {
		name         string
		finder       handlers.URLFinder
		expectedCode codes.Code
		ctx          func() context.Context
	}{
		{
			name:         "нет_пользователя",
			finder:       new(MockFinder),
			expectedCode: codes.InvalidArgument,
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name:         "ошибка_finderа",
			finder:       new(MockFinderBad),
			expectedCode: codes.Internal,
			ctx: func() context.Context {
				md := metadata.New(map[string]string{mData.UserUUID: "1111-2222-3333-444"})
				ctx := metadata.NewOutgoingContext(context.Background(), md)
				return ctx
			},
		},
		{
			name:         "нет_ссылок_для_пользователя",
			finder:       memoryStorage,
			expectedCode: codes.NotFound,
			ctx: func() context.Context {
				md := metadata.New(map[string]string{mData.UserUUID: "1111-2222-3333-444"})
				ctx := metadata.NewOutgoingContext(context.Background(), md)
				return ctx
			},
		},
		{
			name:         "ссылки_есть",
			finder:       memoryStorage,
			expectedCode: codes.OK,
			ctx: func() context.Context {
				_, _ = memoryStorage.CreateUser(models.User{
					UUID: "1111-2222-3333-444",
				})
				id, _ := memoryStorage.Add(models.URL{
					URL:      "http://ya.ru",
					ShortURL: "2ljdsf",
				})
				memoryStorage.LikeURLToUser(id, "1111-2222-3333-444")
				md := metadata.New(map[string]string{mData.UserUUID: "1111-2222-3333-444"})
				ctx := metadata.NewOutgoingContext(context.Background(), md)
				return ctx
			},
		},
	}

	sessionStorage := storage.NewSessionStorage()

	stop := make(chan struct{})
	worker := workers.NewWorker(memoryStorage, stop)

	defer func() {
		stop <- struct{}{}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := grpc.NewServer()
			contract.RegisterUserUrlsHandlerServer(s, NewUserURLsHandler(tt.finder, sessionStorage, worker))
			ctx := tt.ctx()

			dopts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(registerServer(s)),
			}
			conn, err := grpc.NewClient(":///test.server", dopts...)

			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			client := contract.NewUserUrlsHandlerClient(conn)
			r := &empty.Empty{}
			response, err := client.View(ctx, r)

			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.expectedCode {
					t.Error("error code: expected", tt.expectedCode, "received", er.Code())
				}
			}
			if status.Code(err) == codes.OK {
				assert.Equal(t, 2, len(response.Items))
			}
		})
	}
}

func TestUserURLsHandler_Delete(t *testing.T) {
	_ = logger.InitLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	tests := []struct {
		name         string
		ShortUrls    []string
		finder       handlers.URLFinder
		expectedCode codes.Code
		ctx          func() context.Context
	}{
		{
			name:         "нет_пользователя",
			ShortUrls:    []string{},
			finder:       memoryStorage,
			expectedCode: codes.InvalidArgument,
			ctx: func() context.Context {
				return context.Background()
			},
		},

		{
			name:         "ок",
			ShortUrls:    []string{"112s", "dsfsdf"},
			finder:       memoryStorage,
			expectedCode: codes.OK,
			ctx: func() context.Context {
				md := metadata.New(map[string]string{mData.UserUUID: "1111-2222-3333-444"})
				ctx := metadata.NewOutgoingContext(context.Background(), md)
				return ctx
			},
		},
	}

	sessionStorage := storage.NewSessionStorage()

	stop := make(chan struct{})
	worker := workers.NewWorker(memoryStorage, stop)

	defer func() {
		stop <- struct{}{}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := grpc.NewServer()
			contract.RegisterUserUrlsHandlerServer(s, NewUserURLsHandler(tt.finder, sessionStorage, worker))
			ctx := tt.ctx()

			dopts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(registerServer(s)),
			}
			conn, err := grpc.NewClient(":///test.server", dopts...)

			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			client := contract.NewUserUrlsHandlerClient(conn)
			r := &contract.DeleteRequest{
				ShortUrls: tt.ShortUrls,
			}
			_, err = client.Delete(ctx, r)

			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.expectedCode {
					t.Error("error code: expected", tt.expectedCode, "received", er.Code())
				}
			}

		})
	}
}
