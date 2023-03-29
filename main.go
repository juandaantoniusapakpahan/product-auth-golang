package main

import (
	"fmt"
	"product-auth/handlers"
	"product-auth/tokenCookie"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "GGWP")
	})
	e.POST("/user/signup", handlers.SignUpHandler())
	e.POST("/user/signin", handlers.SignInHandler())

	apiGroup := e.Group("/api")
	apiGroup.Use(tokenCookie.TokenRefresherMiddleware)
	apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &tokenCookie.Claims{},
		SigningKey:  []byte(tokenCookie.GetTokenKeySecret()),
		TokenLookup: "cookie:access-token",
	}))
	fmt.Println("Sever started at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
