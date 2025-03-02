package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get port from environment variables with fallback
	port := "8082"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// Initialize router
	router := gin.Default()

	// Add custom template functions
	funcMap := template.FuncMap{
		"title": strings.Title,
	}

	// Set HTML templates with functions
	router.SetFuncMap(funcMap)

	// Load templates and static files
	templatePath := "web/templates/*"
	staticPath := "web/static"

	router.LoadHTMLGlob(templatePath)
	router.Static("/static", staticPath)

	// UI Routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Unixify - Account Management",
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Unixify - Login",
		})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Unixify - Register",
		})
	})

	router.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"title": "Unixify - My Profile",
		})
	})

	// Section routes
	sections := []string{"people", "system", "database", "service"}
	for _, section := range sections {
		section := section // Create a new variable to avoid closure issues
		router.GET("/"+section, func(c *gin.Context) {
			c.HTML(http.StatusOK, "section.html", gin.H{
				"title":   "Unixify - " + strings.Title(section),
				"section": section,
			})
		})
	}

	// Mock API endpoints for authentication
	router.GET("/api/auth/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"isAuthenticated": true,
			"isGuest": true,
			"username": "guest",
			"role": "guest",
		})
	})
	
	// Login endpoint that returns a guest token
	router.POST("/api/auth/login", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		
		// Demo credentials - in a real app this would check against a database
		if credentials.Username == "guest" {
			// Return a guest token
			c.JSON(http.StatusOK, gin.H{
				"token": "guest_token_" + fmt.Sprintf("%d", time.Now().Unix()),
				"user": gin.H{
					"id": 0,
					"username": "guest",
					"email": "guest@example.com",
					"firstName": "Guest",
					"lastName": "User",
					"role": "guest",
					"isGuest": true,
				},
			})
		} else if credentials.Username == "admin" && credentials.Password == "admin" {
			// Return an admin token
			c.JSON(http.StatusOK, gin.H{
				"token": "admin_token_" + fmt.Sprintf("%d", time.Now().Unix()),
				"user": gin.H{
					"id": 1,
					"username": "admin",
					"email": "admin@example.com",
					"firstName": "Admin",
					"lastName": "User",
					"role": "admin",
				},
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
	})

	// Start server
	fmt.Printf("Starting simplified server on port %s...\n", port)
	log.Fatal(router.Run(":" + port))
}