package handlers

import (
	"context"
	"log"
	"testing"

	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestRedirectHandler_Redirect(t *testing.T) {

	type want struct {
		code     codes.Code
		location string
	}
	type request struct {
		id string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "test_#1_короткая_ссылка_преобразуется_в_длинную",
			request: request{
				id: "e98192e19505472476a49f10388428ab",
			},
			want: want{
				code: codes.OK,
			},
		},
		{
			name: "Test_#ссылка_не_передана",
			request: request{
				id: "",
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Test_#3_нет_ссылки",
			request: request{
				id: "123",
			},
			want: want{
				code: codes.NotFound,
			},
		},
	}

	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)

	ctx := context.Background()
	s := grpc.NewServer()
	contract.RegisterRedirectHandlerServer(s, NewRedirectHandler(shortURLService))

	dopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(registerServer(s)),
	}
	conn, err := grpc.NewClient(":///test.server", dopts...)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := contract.NewRedirectHandlerClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := &contract.RedirectRequest{Id: tt.request.id}
			_, err := client.Redirect(ctx, r)

			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.want.code {
					t.Error("error code: expected", tt.want.code, "received", er.Code())
				}
			}
		})
	}

}
