package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "githubEvents/api/docs/githubEvents"
	"githubEvents/api/middlewares"
	"githubEvents/api/routeHandlers"
	"githubEvents/shared"
)

// @title Github Events API
// @version 1.0
// @description This is a github events server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http
func Main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", routeHandlers.HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")
	v1.Use(middlewares.MongoDBMiddleware)

	events := v1.Group("/events")
	events.GET("", routeHandlers.GetAllEventsHandler)
	events.GET("/count", routeHandlers.CountEventsHandler)
	events.GET("/actors", routeHandlers.RecentActorsHandler)
	events.GET("/repos", routeHandlers.RecentReposHandler)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", shared.ApiPort)))
}
