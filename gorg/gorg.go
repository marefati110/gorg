package gorg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/logrusorgru/aurora/v4"
	"github.com/marefati110/gorg/internal/middleware"
	"github.com/marefati110/gorg/internal/utils"
)

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

type Route struct {
	Prefix       string
	Version      string
	AuthRequired bool

	Path    string
	Methods []HttpMethod
	Handler func(e echo.Context) error
	Body    any
	Query   any
	Res     any
}

type ModuleConfig struct {
	Prefix       string
	Version      string
	AuthRequired bool
	Name         string
	Routes       []Route
}

type Config struct {
	Engine                    *echo.Echo
	ModuleConfigs             []ModuleConfig
	Prefix                    string
	Version                   string
	DisableDefaultMiddlewares bool
	DisabledClearTerminal     bool
	ReleaseMode               bool

	HideBanner             bool
	HideProjectInformation bool
	HidePort               bool
}

//

func RegisterModule(...ModuleConfig) (c []ModuleConfig) {

	return c
}

func urlResolve(r Route, m ModuleConfig, cfg Config) string {

	version := cfg.Version
	if m.Version != "" {
		version = m.Version
	}
	if r.Version != "" {
		version = r.Version
	}

	url := cfg.Prefix + m.Prefix + version + fmt.Sprintf("/%s", m.Name) + r.Path

	return url
}

func printInitLog(cfg *Config) error {

	// clear terminal
	if cfg.DisabledClearTerminal {
		if runtime.GOOS == "linux" {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		} else if runtime.GOOS == "windows" {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	}

	wellComeMessage := aurora.Bold("WellCome back!").Blue().String() + " " + aurora.Italic("v0.0.1").String()
	engineInfo := fmt.Sprintf("Engine: %s", aurora.BgBrightCyan(" echo ").String())
	whiteSpace := strings.Repeat(" ", 78-(len(wellComeMessage)+len(engineInfo)))
	fmt.Println(wellComeMessage, whiteSpace, engineInfo)

	content, err := os.ReadFile("./logo.txt")
	if err != nil {
		log.Panic("cannot open logo file")
	}
	fmt.Println(string(content))

	// print information
	fmt.Println(aurora.Bold(" Project information" + strings.Repeat(" ", 31)).BgGray(18))

	if !cfg.ReleaseMode {
		routeCounter := 0

		for _, module := range cfg.ModuleConfigs {
			for _, route := range module.Routes {
				for _, method := range route.Methods {

					url := urlResolve(route, module, *cfg)
					methodS := fmt.Sprintf("%s%s", method, strings.Repeat(" ", 9-len(method)))
					moduleS := fmt.Sprintf("%s%s", aurora.Bold(" "+module.Name+" ").BgGray(8), strings.Repeat(" ", 4))

					handlerName := utils.GetFunctionName(route.Handler)
					functionName := " -→ " + aurora.White(handlerName).String()

					switch method {
					case GET:
						methodS = aurora.Green(methodS).String()
					case POST:
						methodS = aurora.Blue(methodS).String()
					case PUT:
						methodS = aurora.Yellow(methodS).String()
					case DELETE:
						methodS = aurora.Red(methodS).String()
					default:
						methodS = aurora.White(methodS).String()
					}

					fmt.Println(routeCounter+1, methodS, moduleS, url, functionName)

					routeCounter++
				}
			}
		}

	}

	return nil
}

func middlewareFactor(cfg *Config) error {

	e := cfg.Engine

	e.Logger = middleware.GetEchoLogger()
	e.Use(middleware.Hook())

	return nil
}

func routeFactory(cfg *Config) error {

	e := cfg.Engine

	for _, module := range cfg.ModuleConfigs {
		for _, route := range module.Routes {
			for _, method := range route.Methods {

				url := urlResolve(route, module, *cfg)
				e.Add(string(method), url, route.Handler)
			}
		}
	}

	return nil
}

func engineConfig(cfg *Config) error {

	e := cfg.Engine

	e.HideBanner = true
	// e.HidePort = true

	return nil

}

func GorgFactory(cfg *Config) error {

	validate(cfg)

	if err := middlewareFactor(cfg); err != nil {
		return err
	}

	if err := routeFactory(cfg); err != nil {
		return err
	}

	if err := engineConfig(cfg); err != nil {
		return err
	}

	if err := printInitLog(cfg); err != nil {
		return err
	}

	return nil
}
