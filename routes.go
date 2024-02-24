package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/marefati110/gorg/gorg"
)

func Ali(e echo.Context) error {
	return e.String(http.StatusOK, "ok")
}

var MyRoutes = []gorg.GorgRoute{{
	Path:    "/hello",
	Methods: []gorg.HttpMethod{gorg.GET, gorg.POST},
	Handler: Ali,
	Version: "/v1",
}}
