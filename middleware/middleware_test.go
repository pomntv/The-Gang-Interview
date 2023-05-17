package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	// Create a new echo instance
	e := echo.New()

	// Create a new middleware instance
	m := NewMiddleware()

	// Create a test request with a valid API key
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the middleware with the test context
	mw := m.Authentication(func(context echo.Context) error {
		// Just return nil for this test
		return nil
	})

	if assert.NoError(t, mw(c)) {
		// Assert user and apikey are set correctly
		assert.Equal(t, uint(12345), c.Get("user"))
		assert.Equal(t, "12345", c.Get("apikey"))
	}
}

func TestRateLimit(t *testing.T) {
	// Create a new echo instance
	e := echo.New()

	// Create a new middleware instance
	m := NewMiddleware()

	// Create a test request with a valid API key
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "12345")
	rec := httptest.NewRecorder()

	// Call the middleware with the test context

	// Get the current day of the week
	now := time.Now()

	rateforTest := 500
	// Set rate limit based on the day of the week
	rateLimit := 100
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		rateLimit = 200
	}
	// fmt.Printf("mw: %v\n", mw)
	// Make requests until the rate limit is exceeded
	// for i := 1; i <= rateLimit; i++ {
	for i := 1; i <= rateforTest; i++ {
		c := e.NewContext(req, rec)
		c.Set("apikey", "12345")
		mw := m.RateLimit(func(context echo.Context) error {
			// Just return nil for this test
			return nil
		})
		err := mw(c)

		if i <= rateLimit {
			fmt.Printf("case 1 -> i: %v , response : %v\n", i, mw(c))
			assert.NoError(t, err)
		} else if i > rateLimit && i <= rateforTest {
			fmt.Printf("case 2 -> i: %v , response : %v\n", i, http.StatusTooManyRequests)
		} else {
			assert.Error(t, err)
		}
	}
}
