package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)

	if !rl.allow("test-key") {
		t.Error("first request should be allowed")
	}
	if !rl.allow("test-key") {
		t.Error("second request should be allowed")
	}
	if !rl.allow("test-key") {
		t.Error("third request should be allowed")
	}
	if rl.allow("test-key") {
		t.Error("fourth request should be rejected")
	}
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)

	if !rl.allow("key-1") {
		t.Error("key-1 first request should be allowed")
	}
	if !rl.allow("key-2") {
		t.Error("key-2 first request should be allowed")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimit(2, time.Minute))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// First two should succeed
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, httptest.NewRequest("GET", "/test", nil))
	if w1.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w1.Code)
	}

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest("GET", "/test", nil))
	if w2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w2.Code)
	}

	// Third should be rate limited
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, httptest.NewRequest("GET", "/test", nil))
	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w3.Code)
	}
}
