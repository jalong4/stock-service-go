package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jalong4/stock-service-go/auth"
    "github.com/jalong4/stock-service-go/config"
)

type RouteMetadata struct {
    Method       string
    Path         string
    Description  string
    Handler      gin.HandlerFunc
    RequiresAuth bool
}

var routeDefinitions = []RouteMetadata{
    {Method: "POST", Path: "/users/login", Description: "Authenticate user and provide tokens", Handler: LoginHandler, RequiresAuth: false},
    {Method: "POST", Path: "/users/register", Description: "Register a new user", Handler: RegisterUserHandler, RequiresAuth: false},
    {Method: "GET", Path: "/users/", Description: "Retrieve all users", Handler: GetAllUsers, RequiresAuth: true},
    {Method: "GET", Path: "/users/id/:_id", Description: "Retrieve a user by their ID", Handler: GetUserByID, RequiresAuth: true},
    {Method: "DELETE", Path: "/users/id/:_id", Description: "Delete a user by their ID", Handler: DeleteUserHandler, RequiresAuth: true},
    {Method: "PUT", Path: "/users/id/:_id", Description: "Update a user by their ID", Handler: UpdateUserHandler, RequiresAuth: true},
    {Method: "GET", Path: "/holdings/", Description: "Retrieve all holdings", Handler: GetAllHoldingsHandler, RequiresAuth: true},
    {Method: "POST", Path: "/holdings/", Description: "Add a new holding", Handler: AddHoldingsHandler, RequiresAuth: true},
    {Method: "GET", Path: "/holdings/id/:_id", Description: "Retrieve a holding by its ID", Handler: GetHoldingByIDHandler, RequiresAuth: true},
    {Method: "DELETE", Path: "/holdings/id/:_id", Description: "Delete a holding by its ID", Handler: DeleteHoldingHandler, RequiresAuth: true},
    {Method: "PUT", Path: "/holdings/id/:_id", Description: "Update a holding by its ID", Handler: UpdateHoldingHandler, RequiresAuth: true},
    {Method: "GET", Path: "/holdings/ticker/:ticker", Description: "Retrieve holdings by ticker", Handler: GetHoldingsByTickerHandler, RequiresAuth: true},
    {Method: "GET", Path: "/holdings/account/:account", Description: "Retrieve holdings by account", Handler: GetHoldingsByAccountHandler, RequiresAuth: true},
}

func GetRoutes() []RouteMetadata {
    return routeDefinitions
}

func Setup(router *gin.Engine) {
    // Load HTML templates
    router.LoadHTMLGlob("./templates/*.tmpl")

    // Serve static files using dynamically determined base path
    staticBasePath := config.GetBasePath()
    router.Static("/css", staticBasePath + "/css")
    router.Static("/js", staticBasePath + "/js")
    router.Static("/images", staticBasePath + "/images")

    // Define routes dynamically based on routeDefinitions
    for _, route := range routeDefinitions {
        if route.RequiresAuth {
            router.Handle(route.Method, route.Path, auth.AuthMiddleware(), route.Handler)
        } else {
            router.Handle(route.Method, route.Path, route.Handler)
        }
    }

    // Serve the dynamically generated HTML index page at the root
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "Routes": routeDefinitions,
        })
    })
}
