// Package compose provides a helper function to compose Echo middlewares.
package compose

import "github.com/labstack/echo/v4"

// Compose returns a middleware which is composed of given middlewares.
func Compose(middlewares ...echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := next
			for i := len(middlewares) - 1; i >= 0; i-- {
				h = middlewares[i](h)
			}
			return h(c)
		}
	}
}
