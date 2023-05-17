package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type CurrencyAPI struct {
	baseURL string
	cache   *cache.Cache
}

type apiResponse struct {
	Data struct {
		Rates map[string]string `json:"rates"`
	} `json:"data"`
}

func NewCurrencyAPI(baseURL string) *CurrencyAPI {
	return &CurrencyAPI{
		baseURL: baseURL,
		cache:   cache.New(1*time.Hour, 2*time.Hour), // Cache will expire after 1 hour, and cleaned up every 2 hours
	}
}

func (api *CurrencyAPI) FetchExchangeRate(from, to string) (float64, error) {
	// if server have cache
	cacheKey := from + ":" + to
	if rate, found := api.cache.Get(cacheKey); found {
		return rate.(float64), nil
	}
	// no cache use api
	resp, err := http.Get(api.baseURL + "/v2/exchange-rates?currency=" + from)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return 0, err
	}

	rateStr, ok := apiResp.Data.Rates[to]
	if !ok {
		return 0, err
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return 0, err
	}

	api.cache.Set(cacheKey, rate, cache.DefaultExpiration)

	return rate, nil
}
