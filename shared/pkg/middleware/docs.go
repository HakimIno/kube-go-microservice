package middleware

import (
	"context"
	"embed"
	"html/template"

	"github.com/cloudwego/hertz/pkg/app"
)

//go:embed templates/*
var templateFiles embed.FS

// CustomDocsHandler creates a custom API documentation page
func CustomDocsHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		tmpl, err := template.ParseFS(templateFiles, "templates/docs.html")
		if err != nil {
			c.JSON(500, map[string]interface{}{
				"error": "Failed to load documentation template",
			})
			return
		}

		data := map[string]interface{}{
			"Title":       "User Service API Documentation",
			"Description": "Complete API documentation for user management operations",
			"Version":     "1.0.0",
			"BaseURL":     "http://localhost:8081",
			"Endpoints": []map[string]interface{}{
				{
					"Method":      "POST",
					"Path":        "/api/v1/users/register",
					"Description": "Register a new user",
					"Color":       "blue",
				},
				{
					"Method":      "POST",
					"Path":        "/api/v1/users/login",
					"Description": "User login",
					"Color":       "green",
				},
				{
					"Method":      "GET",
					"Path":        "/api/v1/users/{id}",
					"Description": "Get user by ID",
					"Color":       "purple",
				},
				{
					"Method":      "PUT",
					"Path":        "/api/v1/users/{id}",
					"Description": "Update user information",
					"Color":       "orange",
				},
				{
					"Method":      "DELETE",
					"Path":        "/api/v1/users/{id}",
					"Description": "Delete user",
					"Color":       "red",
				},
			},
		}

		c.Header("Content-Type", "text/html")
		tmpl.Execute(c.Response.BodyWriter(), data)
	}
}
