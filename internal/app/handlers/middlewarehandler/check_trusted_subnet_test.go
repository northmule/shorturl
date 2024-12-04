package middlewarehandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"

	"github.com/stretchr/testify/assert"
)

func TestGrantAccess_NoTrustedSubnet(t *testing.T) {
	_ = logger.InitLogger("fatal")
	configApp := &config.Config{
		TrustedSubnet: "",
	}
	middleware := NewCheckTrustedSubnet(configApp)

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.GrantAccess(next).ServeHTTP(res, req)

	assert.Equal(t, http.StatusForbidden, res.Code, fmt.Sprintf("Expected: %d, actual: %d", http.StatusForbidden, res.Code))
}

func TestGrantAccess_InvalidCIDR(t *testing.T) {
	_ = logger.InitLogger("fatal")
	configApp := &config.Config{
		TrustedSubnet: "invalid-cidr",
	}
	middleware := NewCheckTrustedSubnet(configApp)

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.GrantAccess(next).ServeHTTP(res, req)

	assert.Equal(t, http.StatusForbidden, res.Code, fmt.Sprintf("Expected: %d, actual: %d", http.StatusForbidden, res.Code))
}

func TestGrantAccess_IPNotProvided(t *testing.T) {
	_ = logger.InitLogger("fatal")
	configApp := &config.Config{
		TrustedSubnet: "192.168.1.0/24",
	}
	middleware := NewCheckTrustedSubnet(configApp)

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.GrantAccess(next).ServeHTTP(res, req)

	assert.Equal(t, http.StatusForbidden, res.Code, fmt.Sprintf("Expected: %d, actual: %d", http.StatusForbidden, res.Code))
}

func TestGrantAccess_IPNotInSubnet(t *testing.T) {
	_ = logger.InitLogger("fatal")
	configApp := &config.Config{
		TrustedSubnet: "192.168.1.0/24",
	}
	middleware := NewCheckTrustedSubnet(configApp)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")
	res := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.GrantAccess(next).ServeHTTP(res, req)

	assert.Equal(t, http.StatusForbidden, res.Code, fmt.Sprintf("Expected: %d, actual: %d", http.StatusForbidden, res.Code))
}

func TestGrantAccess_IPInSubnet(t *testing.T) {
	_ = logger.InitLogger("fatal")

	configApp := &config.Config{
		TrustedSubnet: "192.168.1.0/24",
	}
	middleware := NewCheckTrustedSubnet(configApp)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.10")
	res := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.GrantAccess(next).ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("Expected: %d, actual: %d", http.StatusOK, res.Code))
}

func TestGrantAccess_IPEqual(t *testing.T) {
	_ = logger.InitLogger("fatal")

	configApp := &config.Config{
		TrustedSubnet: "192.168.1.10/24",
	}
	middleware := NewCheckTrustedSubnet(configApp)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.10")
	res := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.GrantAccess(next).ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("Expected: %d, actual: %d", http.StatusOK, res.Code))
}
