// routes.go
package routes

import (
    "net/http"
    "strings"

    "github.com/jalong4/stock-service-go/auth"
    "github.com/jalong4/stock-service-go/config"
    "github.com/gin-gonic/gin"
)


// Setup initializes all the routes for the application.

func Setup(router *gin.Engine) {
	// Load HTML templates
	router.LoadHTMLGlob("./templates/*.tmpl")

	// Serve static files using dynamically determined base path
	staticBasePath := config.GetBasePath()
	router.Static("/css", staticBasePath + "/css")
	router.Static("/js", staticBasePath + "/js")
	router.Static("/images", staticBasePath + "/images")


    // User routes
	router.POST("/users/login", LoginHandler)
    router.POST("/users/register", RegisterUserHandler)

	router.GET("/users/", auth.AuthMiddleware(), GetAllUsers)
	router.GET("/users/id/:_id", auth.AuthMiddleware(), GetUserByID)
    router.DELETE("/users/id/:_id", auth.AuthMiddleware(), DeleteUserHandler)
    router.PUT("/users/id/:_id", UpdateUserHandler)

	// Holdings routes
	router.GET("/holdings/", auth.AuthMiddleware(), GetAllHoldingsHandler)
    router.POST("/holdings/", auth.AuthMiddleware(), AddHoldingsHandler)
    router.GET("/holdings/id/:_id", auth.AuthMiddleware(), GetHoldingByIDHandler)
    router.DELETE("/holdings/id/:_id", auth.AuthMiddleware(), DeleteHoldingHandler)
    router.PUT("/holdings/id/:_id", auth.AuthMiddleware(), UpdateHoldingHandler)
	router.GET("/holdings/ticker/:ticker", auth.AuthMiddleware(), GetHoldingsByTickerHandler)
	router.GET("/holdings/account/:account", auth.AuthMiddleware(), GetHoldingsByAccountHandler)

	// Serve the dynamically generated HTML index page at the root
	router.GET("/", func(c *gin.Context) {
        allRoutes := router.Routes()
        var apiRoutes []gin.RouteInfo
    
        staticPrefixes := []string{"/css", "/js", "/images"} // static paths
    
        for _, route := range allRoutes {
            includeRoute := true
            for _, prefix := range staticPrefixes {
                if strings.HasPrefix(route.Path, prefix) {
                    includeRoute = false
                    break
                }
            }
            if includeRoute {
                apiRoutes = append(apiRoutes, route)
            }
        }
    
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "Routes": apiRoutes,
        })
	})
}