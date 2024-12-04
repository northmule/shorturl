package utils

import (
	"context"

	mData "github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// FillUserUUID вернёт uuid пользователя
func FillUserUUID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}
	mdValues := md.Get(mData.UserUUID)

	if len(mdValues) == 0 {
		return "", status.Errorf(codes.InvalidArgument, "expected "+mData.UserUUID)
	}

	return mdValues[0], nil
}

// GetUserToken получить токен из запроса.
func GetUserToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	mdValues := md.Get(mData.Authorization)

	if len(mdValues) == 0 {
		return ""
	}

	return mdValues[0]
}

// AppendMData добавит значение в метадату
func AppendMData(ctx context.Context, key string, value string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		md = metadata.New(map[string]string{key: value})
		ctx = metadata.NewIncomingContext(ctx, md)
		return ctx
	}
	mdValues := md.Get(key)
	if len(mdValues) == 0 {
		md.Append(key, value)
		ctx = metadata.NewIncomingContext(ctx, md)
		return ctx
	}
	md.Delete(key)
	md.Append(key, value)
	ctx = metadata.NewIncomingContext(ctx, md)
	return ctx
}
