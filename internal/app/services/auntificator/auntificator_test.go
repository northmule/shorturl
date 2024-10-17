package auntificator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/northmule/shorturl/internal/app/logger"
	"net/http"
	"testing"
	"time"
)

func TestGetUserToken(t *testing.T) {
	_ = logger.NewLogger("info")

	// Test cases
	testCases := []struct {
		name     string
		req      *http.Request
		expected string
	}{
		{
			name: "no_authorization_header_and_no_cookie",
			req: &http.Request{
				Header: http.Header{},
			},
			expected: "",
		},
		{
			name: "authorization_header_present",
			req: &http.Request{
				Header: http.Header{"Authorization": []string{"Bearer token123"}},
			},
			expected: "Bearer token123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetUserToken(tc.req)
			if got != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	_ = logger.NewLogger("info")
	userUUID := "user123"
	exp := time.Hour * 600
	secretKey := "super_secret_key"

	token, expTime := GenerateToken(userUUID, exp, secretKey)

	_, err := hex.DecodeString(token)
	if err != nil {
		t.Errorf("Generated token %s is not a valid hex string", token)
	}

	if expTime.Sub(time.Now()) > exp {
		t.Errorf("Token expiration time is not within the expected range")
	}
}

func TestValidateToken(t *testing.T) {
	_ = logger.NewLogger("info")
	userUUID := "user123"

	// Test cases
	testCases := []struct {
		name     string
		userUUID string
		token    string
		expected bool
	}{
		{
			name:     "invalid_hex_token",
			userUUID: userUUID,
			token:    "invalidToken",
			expected: false,
		},
		{
			name:     "tampered_token",
			userUUID: userUUID,
			token:    hex.EncodeToString(hmac.New(sha256.New, []byte(HMACSecretKey)).Sum([]byte(userUUID))),
			expected: false,
		},
		{
			name:     "valid_token",
			userUUID: userUUID,
			token: func() string {
				token, _ := GenerateToken(userUUID, HMACTokenExp, HMACSecretKey)
				return token
			}(),
			expected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateToken(tc.userUUID, tc.token, HMACSecretKey)
			if got != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	_ = logger.NewLogger("fatal")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateToken("111111-22222-33333-44444", HMACTokenExp, HMACSecretKey)
	}
}

func BenchmarkValidateToken(b *testing.B) {
	_ = logger.NewLogger("fatal")
	userUUID := "user123"
	token, _ := GenerateToken(userUUID, HMACTokenExp, HMACSecretKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateToken(userUUID, token, HMACSecretKey)
	}
}
