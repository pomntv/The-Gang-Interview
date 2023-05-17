package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockCurrencyAPI struct{}

func (m *MockCurrencyAPI) FetchExchangeRate(from, to string) (float64, error) {
	return 50000, nil // Mock exchange rate
}

func TestHandleCurrencyConversion(t *testing.T) {
	// Create a new echo instance
	e := echo.New()

	// Create a new handler with the mock api
	h := NewHandler(&MockCurrencyAPI{})

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/?from=BTC&to=USD&amount=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler with the test context
	if assert.NoError(t, h.HandleCurrencyConversion(c)) {
		// Assert status code and body
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"convertedAmount":"100000.00"}`, strings.TrimRight(rec.Body.String(), "\n"))
	}
}
