package gorg

import "fmt"

func validate(c *Config) {

	e := c.Engine
	if e == nil {
		panic("engine required")
	}

	validateModuleConfig(c)
	validateRouteConfig(c)
}

func validateModuleConfig(c *Config) {
	for _, item := range c.Groups {
		if item.Name == "" {
			panic("Module name required")
		}
	}
}

func validateRouteConfig(c *Config) {
	for _, m := range c.Groups {
		for _, r := range m.Routes {
			if r.Path == "" {
				panic(fmt.Sprintf("module:%s   route path is required", m.Name))
			}

			if r.Handler == nil {
				panic(fmt.Sprintf("module:%s; path:%s  route handler is required", m.Name, r.Path))
			}

			if (r.Methods == nil || len(r.Methods) == 0) && r.Method == "" {
				panic(fmt.Sprintf("module:%s; path:%s  route method is required", m.Name, r.Path))
			}

		}
	}
}
