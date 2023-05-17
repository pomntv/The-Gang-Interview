package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CurrencyAPI is an interface that describes the methods we need for fetching currency rates.
type CurrencyAPI interface {
	FetchExchangeRate(from, to string) (float64, error)
}

type Handler struct {
	CurrencyAPI CurrencyAPI
}

func NewHandler(api CurrencyAPI) *Handler {
	return &Handler{
		CurrencyAPI: api,
	}
}

func (h *Handler) HandleCurrencyConversion(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	amountStr := c.QueryParam("amount")

	// Convert amount to float and handle error
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid amount",
		})
	}

	// Use the api to fetch the exchange rate and handle error
	rate, err := h.CurrencyAPI.FetchExchangeRate(from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch exchange rate",
		})
	}

	// Calculate converted amount
	convertedAmount := amount * rate

	// Respond with the converted amount in JSON format
	return c.JSON(http.StatusOK, map[string]string{
		"convertedAmount": strconv.FormatFloat(convertedAmount, 'f', 2, 64),
	})
}
