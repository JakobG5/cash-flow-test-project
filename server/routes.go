package server

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {

	e.GET("/health", healthCheckHandler)

	v1 := e.Group("/api/v1")
	{
		v1.GET("/", helloHandler)

	}
}
