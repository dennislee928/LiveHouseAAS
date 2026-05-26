package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(nil)

	if !rl.allowInMemory("test-key", 3, time.Minute) {
		t.Error("first request should be allowed")
	}
	if !rl.allowInMemory("test-key", 3, time.Minute) {
		t.Error("second request should be allowed")
	}
	if !rl.allowInMemory("test-key", 3, time.Minute) {
		t.Error("third request should be allowed")
	}
	if rl.allowInMemory("test-key", 3, time.Minute) {
		t.Error("fourth request should be rejected")
	}
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	rl := NewRateLimiter(nil)

	if !rl.allowInMemory("key-1", 1, time.Minute) {
		t.Error("key-1 first request should be allowed")
	}
	if !rl.allowInMemory("key-2", 1, time.Minute) {
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

	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, httptest.NewRequest("GET", "/test", nil))
	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w3.Code)
	}
}
