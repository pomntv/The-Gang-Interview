package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Middleware struct {
	// Can use a database or any other type of storage
	Storage map[string][]time.Time
}

func NewMiddleware() *Middleware {
	return &Middleware{
		Storage: make(map[string][]time.Time),
	}
}

func (m *Middleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		elapsed := time.Since(start)

		log.Printf("[%s] %s, %s %s, status: %d, took: %s\n", start.Format(time.RFC3339), c.Request().RemoteAddr, c.Request().Method, c.Request().URL, c.Response().Status, elapsed)

		return err
	}
}

func (m *Middleware) Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered: %v\n%s", r, debug.Stack())
				c.NoContent(http.StatusInternalServerError)
			}
		}()
		return next(c)
	}
}

func (m *Middleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get API key from the request header
		apiKey := c.Request().Header.Get("X-API-Key")

		// Validate API key and get user ID
		userID, err := m.ValidateAPIKey(apiKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key")
		}

		// Store user ID,apikey in the request context
		c.Set("user", userID)
		c.Set("apikey", apiKey)
		return next(c)
	}
}

func (m *Middleware) ValidateAPIKey(apiKey string) (uint, error) { // Fake mapping API and ID API == ID
	userID64, err := strconv.ParseInt(apiKey, 10, 32)
	if err != nil {
		return 0, err
	}
	userID := uint(userID64)
	return userID, nil
}

func (m *Middleware) RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// fmt.Printf("c: %T\n", c)  // check c
		// fmt.Printf("%+v\n", c)
		apiKey, ok := c.Get("apikey").(string)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "missing API key")
		}

		// Check the rate limit for the user associated with the API key.
		requestTimes, ok := m.Storage[apiKey]
		if !ok {
			requestTimes = []time.Time{}
		}

		// Check if the user has exceeded their rate limit
		now := time.Now()
		var newRequestTimes []time.Time
		for _, requestTime := range requestTimes {
			if now.Sub(requestTime) < 24*time.Hour {
				newRequestTimes = append(newRequestTimes, requestTime)
			}
		}
		// newRequestTimes = append(newRequestTimes, now)

		// The rate limit is 100 requests per workday (Monday-Friday) and 200 requests per day on weekends
		rateLimit := 100
		checkprintrateLimit := "rate limit exceeded (workday rateLimit 100)"
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			rateLimit = 200
			checkprintrateLimit = "rate limit exceeded (weekday rateLimit 200)"
		}
		if len(newRequestTimes) > rateLimit {
			return echo.NewHTTPError(http.StatusTooManyRequests, checkprintrateLimit) // 429 too Many Requests
		}

		m.Storage[apiKey] = newRequestTimes

		return next(c)
	}
}
