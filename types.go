package gorg

import "github.com/labstack/echo/v4"

type HttpMethod string
type GorgMiddleware string

const (
	GET     HttpMethod = "GET"
	PUT     HttpMethod = "PUT"
	HEAD    HttpMethod = "HEAD"
	POST    HttpMethod = "POST"
	PATCH   HttpMethod = "PATCH"
	TRACE   HttpMethod = "TRACE"
	DELETE  HttpMethod = "DELETE"
	OPTIONS HttpMethod = "OPTIONS"
	CONNECT HttpMethod = "CONNECT"
)

const (
	Auth    GorgMiddleware = "auth"
	Logger  GorgMiddleware = "logger"
	Recover GorgMiddleware = "recover"
)

type RouteDoc struct {
	Description string
	Summary     string
}

type Route struct {
	Prefix       string
	Version      string
	AuthRequired bool

	Doc RouteDoc

	Path    string
	Method  HttpMethod
	Methods []HttpMethod
	Handler func(c echo.Context) error
	Body    any
	Query   any
	Res     any
}

type Group struct {
	Prefix       string
	Version      string
	AuthRequired bool
	Name         string
	Routes       []Route
	Middlewares  []echo.MiddlewareFunc
}

type Swagger struct {
	IsDisabled  bool
	ID          string
	Swagger     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

type Config struct {
	Swagger

	Engine *echo.Echo
	Groups []Group

	Prefix  string
	Version string

	// DisableDefaultMiddlewares bool
	DisabledClearTerminal bool
	ReleaseMode           bool

	HideBanner      bool
	HidePort        bool
	HideInformation bool
}
