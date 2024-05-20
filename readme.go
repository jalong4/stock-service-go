package main

import (
    "log"
    "os"
    "text/template"

    "github.com/gin-gonic/gin"
    "github.com/jalong4/stock-service-go/routes"
)

type RouteInfo struct {
    Method       string
    Path         string
    Description  string
    RequiresAuth bool
}

type ReadmeData struct {
    Routes []RouteInfo
}

func GenerateReadme(router *gin.Engine) {
    // Collect routes information using GetRoutes
    ginRoutes := routes.GetRoutes()
    var routeInfo []RouteInfo
    for _, route := range ginRoutes {
        routeInfo = append(routeInfo, RouteInfo{
            Method:       route.Method,
            Path:         route.Path,
            Description:  route.Description,
            RequiresAuth: route.RequiresAuth,
        })
    }

    data := ReadmeData{
        Routes: routeInfo,
    }

    // Parse the README template
    tmpl, err := template.ParseFiles("templates/README.tmpl")
    if err != nil {
        log.Fatalf("Failed to parse README template: %v", err)
    }

    // Create the README.md file
    file, err := os.Create("README.md")
    if err != nil {
        log.Fatalf("Failed to create README.md: %v", err)
    }
    defer file.Close()

    // Execute the template and write to the README.md file
    err = tmpl.Execute(file, data)
    if err != nil {
        log.Fatalf("Failed to generate README.md: %v", err)
    }

    log.Println("README.md generated successfully")
}
