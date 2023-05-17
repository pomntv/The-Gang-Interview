package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchExchangeRate(t *testing.T) {
	// Create a test server to mock the API response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"data": {
				"rates": {
					"USD": "50000"
				}
			}
		}`))
	}))
	defer ts.Close()

	// Replace baseURL with our test server address
	api := NewCurrencyAPI(ts.URL)

	rate, err := api.FetchExchangeRate("BTC", "USD")
	if err != nil {
		t.Errorf("Failed to fetch exchange rate: %s", err)
	}

	// Check if the returned rate is as expected
	if rate != 50000 {
		t.Errorf("Expected rate to be 50000, but got %f", rate)
	}

	// Test caching by fetching the rate again
	rate, err = api.FetchExchangeRate("BTC", "USD")
	if err != nil {
		t.Errorf("Failed to fetch exchange rate from cache: %s", err)
	}

	// Check if the returned rate from cache is as expected
	if rate != 50000 {
		t.Errorf("Expected cached rate to be 50000, but got %f", rate)
	}
}
