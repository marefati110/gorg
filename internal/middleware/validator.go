package middleware

import (
	"github.com/labstack/echo/v4"
)

func Validator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
