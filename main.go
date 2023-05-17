package main

import (
	"The-Gang-Interview/api"
	"The-Gang-Interview/handlers"
	"The-Gang-Interview/middleware"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware" // Logger Recover build in is appear in my middleware.go
)

func main() {
	e := echo.New()

	//  New Instance Middelware
	m := middleware.NewMiddleware()
	e.Use(m.Logger)
	e.Use(m.Recover)
	e.Use(m.Authentication)
	e.Use(m.RateLimit)

	// Routes
	api := api.NewCurrencyAPI("https://api.coinbase.com")

	// Create a Handler with the api
	handler := handlers.NewHandler(api)
	e.GET("/convert", handler.HandleCurrencyConversion)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

}
