package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ApplyPandaContext)
	e.Use(RetrieveCredentials)

	e.Any("/*", handle)

	e.Logger.Fatal(e.Start(":60394"))
}
