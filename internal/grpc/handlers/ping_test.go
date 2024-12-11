package handlers

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type MockPostgresStorageOk struct {
	mock.Mock
}

func (m *MockPostgresStorageOk) Add(url models.URL) (int64, error) {
	return 0, nil
}
func (m *MockPostgresStorageOk) FindByShortURL(shortURL string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageOk) FindByURL(url string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageOk) Ping() error {
	return nil
}
func (m *MockPostgresStorageOk) MultiAdd(url []models.URL) error {
	return nil
}
func (m *MockPostgresStorageOk) CreateUser(user models.User) (int64, error) {
	return 0, nil
}

func (m *MockPostgresStorageOk) LikeURLToUser(urlID int64, userUUID string) error {
	return nil
}

func (m *MockPostgresStorageOk) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}

func (m *MockPostgresStorageOk) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	return nil
}

type MockPostgresStorageBad struct {
	mock.Mock
}

func (m *MockPostgresStorageBad) Add(url models.URL) (int64, error) {
	return 0, nil
}
func (m *MockPostgresStorageBad) FindByShortURL(shortURL string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageBad) FindByURL(url string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageBad) Ping() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockPostgresStorageBad) MultiAdd(url []models.URL) error {
	return nil
}
func (m *MockPostgresStorageBad) CreateUser(user models.User) (int64, error) {
	return 0, nil
}

func (m *MockPostgresStorageBad) LikeURLToUser(urlID int64, userUUID string) error {
	return nil
}

func (m *MockPostgresStorageBad) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}

func (m *MockPostgresStorageBad) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	return nil
}

func registerServer(s *grpc.Server) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestPingHandler_CheckStorageConnect(t *testing.T) {

	file, err := os.CreateTemp("/tmp", "TestFileStorage_Add_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	postgresStorage := new(MockPostgresStorageOk)
	postgresStorage.On("Ping").Return("Ok")

	tests := []struct {
		name    string
		storage storage.StorageQuery
		want    bool
	}{
		{
			name:    "MemoryStorage",
			storage: storage.NewMemoryStorage(),
			want:    true,
		},
		{
			name:    "FileStorage",
			storage: storage.NewFileStorage(file),
			want:    true,
		},
		{
			name:    "PostgresStorage",
			storage: postgresStorage,
			want:    true,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := grpc.NewServer()
			contract.RegisterPingHandlerServer(s, NewPingHandler(tt.storage))

			dopts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(registerServer(s)),
			}
			conn, err := grpc.NewClient(":///test.server", dopts...)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := contract.NewPingHandlerClient(conn)
			request := &empty.Empty{}
			response, err := client.CheckStorageConnect(ctx, request)

			if err != nil {
				t.Fatal(err)
			}

			if response.Ok != tt.want {
				t.Errorf("want %v; got %v", tt.want, response.Ok)
			}
		})
	}

	t.Run("Возврат_ошибки_подключения", func(t *testing.T) {
		mockStorage := new(MockPostgresStorageBad)
		mockStorage.On("Ping").Return(errors.New("bad test request"))
		s := grpc.NewServer()
		contract.RegisterPingHandlerServer(s, NewPingHandler(mockStorage))

		dopts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(registerServer(s)),
		}
		conn, err := grpc.NewClient(":///test.server", dopts...)

		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := contract.NewPingHandlerClient(conn)
		request := &empty.Empty{}
		response, err := client.CheckStorageConnect(ctx, request)

		if err == nil {
			t.Error("expected error")
		}
		if er, ok := status.FromError(err); ok {
			if er.Code() != codes.Internal {
				t.Error("error code: expected", codes.Internal, "received", er.Code())
			}
			if er.Message() != "no connect db" {
				t.Error("error message: expected", "no connect db", "received", er.Message())
			}
		}
		if response != nil {
			t.Error("not expected response")
		}
	})
}
