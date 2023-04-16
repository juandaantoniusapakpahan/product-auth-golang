package main

import (
	"fmt"
	"net/http"
	"product-auth/handlers"
	"product-auth/tokenCookie"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4040"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderAccept, echo.HeaderContentType},
		AllowCredentials: true,
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
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
	apiGroup.POST("/resume", handlers.CreateResumHandler())
	apiGroup.PUT("/resume/:id", handlers.EditResumeHandler())
	apiGroup.GET("/resume", handlers.GetResumeByPrivateHandler())
	apiGroup.DELETE("/resume/:id", handlers.DeleteResumeHandler())

	fmt.Println("Sever started at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
