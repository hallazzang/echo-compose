# echo-compose

[![Godoc Reference](https://godoc.org/github.com/hallazzang/echo-compose?status.svg)](https://godoc.org/github.com/hallazzang/echo-compose)
[![Goreportcard Result](https://goreportcard.com/badge/github.com/hallazzang/echo-compose)](https://goreportcard.com/report/github.com/hallazzang/echo-compose)

Compose two or more echo middlewares together.

## Installation

```
go get github.com/hallazzang/echo-compose
```

## Usage Example

```go
package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	compose "github.com/hallazzang/echo-compose"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type customClaims struct {
	Key string `json:"key"`
	jwt.StandardClaims
}

// setClaims extracts JWT token from context and puts parsed claims
// inside the context under "claims" key.
func setClaims(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*customClaims)
		c.Set("claims", claims)
		return next(c)
	}
}

func main() {
	e := echo.New()

	jwtKey := []byte("secret")

	// compose echo JWT middleware with our custom middleware into
	// single middleware.
	m := compose.Compose(
		middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey: jwtKey,
			Claims:     &customClaims{},
		}),
		setClaims,
	)

	e.POST("/token", func(c echo.Context) error {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
			Key: "magical key",
		})
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
	})

	// wrap "/key" endpoint with our composed middleware so that
	// our handler could access to parsed claims.
	e.GET("/key", func(c echo.Context) error {
		claims := c.Get("claims").(*customClaims)
		return c.JSON(http.StatusOK, echo.Map{"key": claims.Key})
	}, m)

	e.Logger.Fatal(e.Start(":5000"))
}
```

You can use your favorite HTTP client to test our server:

```
$ http POST :5000/token
HTTP/1.1 200 OK
Content-Length: 122
Content-Type: application/json; charset=UTF-8
Date: Wed, 01 Apr 2020 11:31:47 GMT

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJtYWdpY2FsIGtleSJ9.mxBhplkeaT3OskFbD_G8xtQ-7uMXDzEB8J5OktIbzUc"
}

$ http :5000/key Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJtYWdpY2FsIGtleSJ9.mxBhplkeaT3OskFbD_G8xtQ-7uMXDzEB8J5OktIbzUc'
HTTP/1.1 200 OK
Content-Length: 22
Content-Type: application/json; charset=UTF-8
Date: Wed, 01 Apr 2020 11:32:37 GMT

{
    "key": "magical key"
}
```
