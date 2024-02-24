package gorg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/logrusorgru/aurora/v4"
	"github.com/marefati110/gorg/internal/middleware"
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

type GorgRoute struct {
	Path    string
	Methods []HttpMethod
	Handler func(e echo.Context) error
	Version string

	Body  any
	Query any
	Res   any
}

type ModuleConfig struct {
	Name string
}

type GorgConfig struct {
	Engine                    *echo.Echo
	Routes                    []GorgRoute
	Prefix                    string
	Version                   string
	DisableDefaultMiddlewares bool
	DisabledClearTerminal     bool
	ReleaseMode               bool
}

//

func urlResolve(r GorgRoute, prefix, defaultVersion string) string {

	version := defaultVersion
	if r.Version == "" {
		version = r.Version
	}

	return prefix + version + r.Path
}

func printInitLog(cfg *GorgConfig) error {

	// clear terminal
	if !cfg.DisabledClearTerminal {
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
		for _, route := range cfg.Routes {

			for _, method := range route.Methods {

				url := urlResolve(route, cfg.Prefix, cfg.Version)

				methodS := fmt.Sprintf("%s%s", method, strings.Repeat(" ", 8-len(method)))

				handlerValue := reflect.ValueOf(route.Handler)
				handlerName := runtime.FuncForPC(handlerValue.Pointer()).Name()
				lastDotIndex := strings.LastIndex(handlerName, ".")
				functionName := " ==> " + "(" + aurora.Bold(handlerName[lastDotIndex+1:]).String() + ")"

				if method == GET {
					fmt.Println(routeCounter+1, aurora.Green(methodS), url, functionName)
				} else if method == POST {
					fmt.Println(routeCounter+1, aurora.Blue(methodS), url)
				} else if method == PUT {
					fmt.Println(routeCounter+1, aurora.Yellow(methodS), url)
				} else if method == DELETE {
					fmt.Println(routeCounter+1, aurora.Red(methodS), url)
				} else {
					fmt.Println(routeCounter+1, aurora.White(methodS), url)
				}

				routeCounter++
			}

		}
	}

	return nil
}

func middlewareFactor(cfg *GorgConfig) error {

	e := cfg.Engine

	e.Logger = middleware.GetEchoLogger()
	e.Use(middleware.Hook())

	return nil
}

func routeFactory(cfg *GorgConfig) error {

	routes := cfg.Routes
	e := cfg.Engine

	prefix := cfg.Prefix

	fmt.Println("hello")

	for _, route := range routes {

		for _, method := range route.Methods {

			url := urlResolve(route, prefix, cfg.Version)
			e.Add(string(method), url, route.Handler)
		}
	}

	return nil
}

func engineConfig(cfg *GorgConfig) error {

	e := cfg.Engine

	e.HideBanner = true

	return nil

}

func GorgFactory(cfg *GorgConfig) error {

	e := cfg.Engine
	if e == nil {
		return fmt.Errorf("engine required")
	}

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
