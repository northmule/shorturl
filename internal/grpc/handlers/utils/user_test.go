package utils

import (
	"context"
	"testing"

	mData "github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestFillUserUUID_MetadataMissing(t *testing.T) {
	ctx := context.Background()
	userUUID, err := FillUserUUID(ctx)

	assert.Equal(t, "", userUUID)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestFillUserUUID_UserUUIDMissing(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs())
	userUUID, err := FillUserUUID(ctx)

	assert.Equal(t, "", userUUID)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestFillUserUUID_UserUUIDPresent(t *testing.T) {
	expectedUUID := "test-uuid"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(mData.UserUUID, expectedUUID))
	userUUID, err := FillUserUUID(ctx)

	assert.Equal(t, expectedUUID, userUUID)
	assert.Nil(t, err)
}

func TestGetUserToken_MetadataMissing(t *testing.T) {
	ctx := context.Background()
	token := GetUserToken(ctx)

	assert.Equal(t, "", token)
}

func TestGetUserToken_AuthorizationMissing(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs())
	token := GetUserToken(ctx)

	assert.Equal(t, "", token)
}

func TestGetUserToken_AuthorizationPresent(t *testing.T) {
	expectedToken := "test-token"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(mData.Authorization, expectedToken))
	token := GetUserToken(ctx)

	assert.Equal(t, expectedToken, token)
}

func TestAppendMData_MetadataMissing(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	newCtx := AppendMData(ctx, key, value)

	md, ok := metadata.FromIncomingContext(newCtx)
	assert.True(t, ok)
	assert.Equal(t, []string{value}, md.Get(key))
}

func TestAppendMData_MetadataPresent_KeyMissing(t *testing.T) {
	key := "test-key"
	value := "test-value"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("other-key", "other-value"))
	newCtx := AppendMData(ctx, key, value)

	md, ok := metadata.FromIncomingContext(newCtx)
	assert.True(t, ok)
	assert.Equal(t, []string{value}, md.Get(key))

}
func TestAppendMData_MetadataPresent_KeyPresent(t *testing.T) {
	key := "test-key"
	value := "test-value"
	newValue := "new-test-value"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(key, value))
	newCtx := AppendMData(ctx, key, newValue)

	md, ok := metadata.FromIncomingContext(newCtx)
	assert.True(t, ok)
	assert.Equal(t, []string{newValue}, md.Get(key))
}

func BenchmarkFillUserUUID(b *testing.B) {
	md := metadata.Pairs(mData.UserUUID, "test-uuid")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	for i := 0; i < b.N; i++ {
		_, _ = FillUserUUID(ctx)
	}
}

func BenchmarkGetUserToken(b *testing.B) {
	md := metadata.Pairs(mData.Authorization, "test-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	for i := 0; i < b.N; i++ {
		_ = GetUserToken(ctx)
	}
}

func BenchmarkAppendMData_NewMetadata(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_ = AppendMData(ctx, "test-key", "test-value")
	}
}

func BenchmarkAppendMData_ExistingMetadata(b *testing.B) {
	md := metadata.Pairs("test-key", "old-value")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	for i := 0; i < b.N; i++ {
		_ = AppendMData(ctx, "test-key", "test-value")
	}
}
