package handlers

import (
	"context"
	"log"
	"testing"

	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/grpc/contract"
	mData "github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestShortenerHandler_Shortener(t *testing.T) {

	tests := []struct {
		name string
		url  string
		code codes.Code
		ctx  func() context.Context
	}{
		{
			name: "ссылка_не_валидная",
			url:  "Жил был слон!",
			code: codes.InvalidArgument,
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "нет_пользователя",
			url:  "https://ya.ru/map1",
			code: codes.InvalidArgument,
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "короткая_ссылка_создаётся",
			url:  "https://ya.ru/map1",
			code: codes.OK,
			ctx: func() context.Context {
				md := metadata.New(map[string]string{mData.UserUUID: "1111-2222-3333-444"})
				ctx := metadata.NewOutgoingContext(context.Background(), md)
				return ctx
			},
		},
	}

	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)

	s := grpc.NewServer()
	contract.RegisterShortenerHandlerServer(s, NewShortenerHandler(shortURLService, memoryStorage, memoryStorage))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
			client := contract.NewShortenerHandlerClient(conn)

			r := &contract.ShortenerRequest{Url: tt.url}
			response, err := client.Shortener(ctx, r)

			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.code {
					t.Error("error code: expected", tt.code, "received", er.Code())
					assert.Equal(t, "", response.ShortUrl)
				}
			}
			if status.Code(err) == codes.OK {
				assert.NotEqual(t, "", response.ShortUrl)
			}
		})
	}
}

func TestShortenerHandler_ShortenerJSON(t *testing.T) {

	tests := []struct {
		name string
		url  string
		code codes.Code
		ctx  func() context.Context
	}{
		{
			name: "ссылка_не_валидная",
			url:  "Жил был слон!",
			code: codes.InvalidArgument,
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "нет_пользователя",
			url:  "https://ya.ru/map1",
			code: codes.InvalidArgument,
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "короткая_ссылка_создаётся",
			url:  "https://ya.ru/map1",
			code: codes.OK,
			ctx: func() context.Context {
				md := metadata.New(map[string]string{mData.UserUUID: "1111-2222-3333-444"})
				ctx := metadata.NewOutgoingContext(context.Background(), md)
				return ctx
			},
		},
	}

	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)

	s := grpc.NewServer()
	contract.RegisterShortenerHandlerServer(s, NewShortenerHandler(shortURLService, memoryStorage, memoryStorage))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
			client := contract.NewShortenerHandlerClient(conn)

			r := &contract.ShortenerJSONRequest{Url: tt.url}
			response, err := client.ShortenerJSON(ctx, r)

			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.code {
					t.Error("error code: expected", tt.code, "received", er.Code())
					assert.Equal(t, "", response.Result)
				}
			}
			if status.Code(err) == codes.OK {
				assert.NotEqual(t, "", response.Result)
			}
		})
	}
}

func TestShortenerHandler_ShortenerBatch(t *testing.T) {

	tests := []struct {
		name  string
		items []*contract.ShortenerBatchRequest_Item
		code  codes.Code
	}{
		{
			name:  "не_переданы_items",
			items: []*contract.ShortenerBatchRequest_Item{},
			code:  codes.InvalidArgument,
		},
		{
			name: "items_и_не_валидные_ссылки",
			items: []*contract.ShortenerBatchRequest_Item{{
				CorrelationId: "1",
				OriginalUrl:   "Парам_пам_пам",
			}},
			code: codes.InvalidArgument,
		},
		{
			name: "ссылки_валидные",
			items: []*contract.ShortenerBatchRequest_Item{{
				CorrelationId: "1",
				OriginalUrl:   "http://ya.ru/result",
			}},
			code: codes.OK,
		},
	}

	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)

	s := grpc.NewServer()
	contract.RegisterShortenerHandlerServer(s, NewShortenerHandler(shortURLService, memoryStorage, memoryStorage))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			dopts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(registerServer(s)),
			}
			conn, err := grpc.NewClient(":///test.server", dopts...)

			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			client := contract.NewShortenerHandlerClient(conn)

			r := &contract.ShortenerBatchRequest{Items: tt.items}
			response, err := client.ShortenerBatch(ctx, r)

			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.code {
					t.Error("error code: expected", tt.code, "received", er.Code())
				}
			}

			if status.Code(err) == codes.OK {
				assert.Equal(t, 1, len(response.Items))
			}
		})
	}
}
