package compose

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func handler(c echo.Context) error {
	return c.String(http.StatusOK, c.Get("secret").(string))
}

func createMiddleware(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			secret, _ := c.Get("secret").(string)
			c.Set("secret", secret+name)
			return next(c)
		}
	}
}

func TestCompose(t *testing.T) {
	e := echo.New()
	m := Compose(
		createMiddleware("a"),
		createMiddleware("b"),
		createMiddleware("c"),
	)
	e.Use(m)
	e.GET("/secret", handler)

	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, rec.Body.String(), "abc")
}
