package gorg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-openapi/spec"
	echojwt "github.com/labstack/echo-jwt/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/logrusorgru/aurora/v4"
	"github.com/marefati110/gorg/internal/middleware"
	"github.com/marefati110/gorg/internal/swagger"
	"github.com/marefati110/gorg/internal/utils"
)

func printLogo() {
	fmt.Println(`
###################################################
#     _______   ______   .______       _______    #
#    /  _____| /  __  \  |   _  \     /  _____|   #
#   |  |  __  |  |  |  | |  |_)  |   |  |  __     #
#   |  | |_ | |  |  |  | |      /    |  | |_ |    #
#   |  |__| | |  |--|  | |  |\  \----|  |__| |    #
#    \______|  \______/  | _| \._____|\______|    #
#                                                 #
###################################################`)
}

func RegisterGroup(modules ...Group) []Group {
	moduleMap := make(map[string]Group)

	for _, module := range modules {
		_, exists := moduleMap[module.Name]
		if !exists {
			moduleMap[module.Name] = module
		}
	}

	uniqueModules := make([]Group, 0, len(moduleMap))
	for _, module := range moduleMap {
		uniqueModules = append(uniqueModules, module)
	}

	return uniqueModules
}

func UrlResolve(r Route, m Group, c Config) string {

	version := c.Version
	if m.Version != "" {
		version = m.Version
	}
	if r.Version != "" {
		version = r.Version
	}

	prefix := m.Prefix
	if r.Prefix != "" {
		prefix = r.Prefix
	}

	url := c.Prefix + prefix + version + fmt.Sprintf("/%s", m.Name) + r.Path

	return url
}

func MethodResolve(r Route) []HttpMethod {

	if r.Method != "" {
		r.Methods = append(r.Methods, r.Method)
	}

	return r.Methods
}

func printInitLog(c *Config) error {

	// clear terminal
	if c.DisabledClearTerminal {
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

	printLogo()

	// print information
	fmt.Println(aurora.Bold(" Project information" + strings.Repeat(" ", 31)).BgGray(18))

	if !c.ReleaseMode {
		routeCounter := 0

		for _, module := range c.Groups {
			for _, route := range module.Routes {

				methods := route.Methods

				if route.Method != "" {
					methods = append(methods, route.Method)

				}

				for _, method := range methods {

					url := UrlResolve(route, module, *c)
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

func middlewareFactor(c *Config) error {

	e := c.Engine

	e.Logger = middleware.GetEchoLogger()
	e.Use(middleware.Hook(c.ReleaseMode))

	e.Use(echoMiddleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
	}))

	e.Use(middleware.Validator())

	return nil
}

func routeFactory(c *Config) error {

	e := c.Engine

	for _, group := range c.Groups {

		for _, r := range group.Routes {

			methods := r.Methods

			if r.Method != "" {
				methods = append(methods, r.Method)
			}

			for _, m := range methods {

				url := UrlResolve(r, group, *c)

				if group.Name == "" {
					g := e.Group(group.Name)
					g.Add(string(m), url, r.Handler, group.Middlewares...)
				} else {
					e.Add(string(m), url, r.Handler, group.Middlewares...)
				}

			}
		}
	}

	return nil
}

func engineConfig(c *Config) error {

	e := c.Engine

	e.HideBanner = true
	// e.HidePort = true

	return nil
}

func swaggerFactor(c *Config) error {

	s := swagger.SwaggerConfig{}
	s.Init()
	s.SetVersion()
	s.SetInfo(spec.InfoProps{
		Title:       c.Swagger.Title,
		Description: c.Swagger.Description,
	})

	return nil
}

func GorgFactory(c *Config) error {

	validate(c)

	if err := middlewareFactor(c); err != nil {
		return err
	}

	if err := routeFactory(c); err != nil {
		return err
	}

	if err := engineConfig(c); err != nil {
		return err
	}

	if err := swaggerFactor(c); err != nil {
		return err
	}

	if err := printInitLog(c); err != nil {
		return err
	}

	return nil
}
