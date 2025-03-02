package api

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/home/unixify/internal/auth"
	"github.com/home/unixify/internal/config"
	"github.com/home/unixify/internal/handlers"
	"github.com/home/unixify/internal/repository"
	"github.com/home/unixify/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Server represents the API server
type Server struct {
	router  *gin.Engine
	config  *config.Config
	logger  *logrus.Logger
	handler *handlers.Handler
	db      *gorm.DB
	repo    *repository.Repository
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, services *service.Services) *Server {
	// Initialize logger
	logger := logrus.New()
	if cfg.Server.Mode == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize router
	router := gin.New()
	
	// Use middleware
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))
	
	// Initialize handlers
	handler := handlers.NewHandler(services, logger)

	// Get database connection and repository
	db := services.GetDB()
	repo := repository.NewRepository(db)

	// Create server
	server := &Server{
		router:  router,
		config:  cfg,
		logger:  logger,
		handler: handler,
		db:      db,
		repo:    repo,
	}

	// Initialize routes
	server.initRoutes()

	return server
}

// LoggerMiddleware is a middleware function that logs requests using logrus
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start).Milliseconds()

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		
		if raw != "" {
			path = path + "?" + raw
		}

		entry := logger.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency_ms": latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
		})

		if statusCode >= 500 {
			entry.Error("Server error")
		} else if statusCode >= 400 {
			entry.Warn("Client error")
		} else {
			entry.Info("Request")
		}
	}
}

// initRoutes initializes the API routes
func (s *Server) initRoutes() {
	// Authentication middleware
	authService := auth.NewService(*s.config)
	authMiddleware := authService.AuthMiddleware()
	guestMiddleware := authService.GuestMiddleware()

	// API routes
	api := s.router.Group("/api")
	{
		// Health check (public)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Guest accessible routes - read-only operations
		guestAPI := api.Group("/")
		guestAPI.Use(guestMiddleware)
		{
			// Account read-only routes
			accounts := guestAPI.Group("/accounts")
			{
				accounts.GET("", s.handler.GetAllAccounts)
				accounts.GET("/:id", s.handler.GetAccount)
				accounts.GET("/uid/:uid", s.handler.GetAccountByUID)
				accounts.GET("/username/:username", s.handler.GetAccountByUsername)
				accounts.GET("/:id/groups", s.handler.GetAccountGroups)
			}

			// Group read-only routes
			groups := guestAPI.Group("/groups")
			{
				groups.GET("", s.handler.GetAllGroups)
				groups.GET("/:id", s.handler.GetGroup)
				groups.GET("/gid/:gid", s.handler.GetGroupByGID)
				groups.GET("/groupname/:groupname", s.handler.GetGroupByGroupname)
				groups.GET("/:id/accounts", s.handler.GetGroupMembers)
			}

			// Search routes (read-only)
			search := guestAPI.Group("/search")
			{
				search.GET("/accounts", s.handler.SearchAccounts)
				search.GET("/groups", s.handler.SearchGroups)
			}

			// Audit routes (read-only)
			audit := guestAPI.Group("/audit")
			{
				audit.GET("", s.handler.GetAuditEntries)
				audit.GET("/:id", s.handler.GetAuditEntry)
			}
		}

		// Protected API routes - require authentication for write operations
		protected := api.Group("/")
		protected.Use(authMiddleware)
		{
			// Account write operations
			accounts := protected.Group("/accounts")
			{
				accounts.POST("", s.handler.CreateAccount)
				accounts.PUT("/:id", s.handler.UpdateAccount)
				accounts.DELETE("/:id", s.handler.DeleteAccount)
			}

			// Group write operations
			groups := protected.Group("/groups")
			{
				groups.POST("", s.handler.CreateGroup)
				groups.PUT("/:id", s.handler.UpdateGroup)
				groups.DELETE("/:id", s.handler.DeleteGroup)
			}

			// Membership routes (write operations)
			membership := protected.Group("/memberships")
			{
				membership.POST("", s.handler.AssignAccountToGroup)
				membership.DELETE("", s.handler.RemoveAccountFromGroup)
			}
		}
	}

	// Auth handler instance
	authHandler := handlers.NewAuthHandler(s.db, auth.NewService(*s.config), s.repo)
	
	// Auth routes
	authRoutes := s.router.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/verify-totp", authHandler.VerifyTOTP)
		
		// Protected auth routes
		protected := authRoutes.Group("/")
		protected.Use(authMiddleware)
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.POST("/update-password", authHandler.UpdatePassword)
			protected.GET("/setup-totp", authHandler.SetupTOTP)
			protected.POST("/activate-totp", authHandler.ActivateTOTP)
			protected.POST("/disable-totp", authHandler.DisableTOTP)
		}
	}

	// Web interface routes
	// Add custom template functions
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	
	// Set HTML templates with functions
	s.router.SetFuncMap(funcMap)
	
	// Get template and static paths from environment variables with fallbacks
	templatePath := "web/templates/*"
	if envPath := s.config.GetString("TEMPLATE_PATH"); envPath != "" {
		templatePath = envPath
	}
	s.logger.Infof("Loading templates from: %s", templatePath)
	
	staticPath := "web/static"
	if envPath := s.config.GetString("STATIC_PATH"); envPath != "" {
		staticPath = envPath
	}
	s.logger.Infof("Serving static files from: %s", staticPath)
	
	s.router.LoadHTMLGlob(templatePath)
	s.router.Static("/static", staticPath)
	
	// Login/profile routes
	s.router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Unixify - Login",
		})
	})
	
	s.router.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"title": "Unixify - My Profile",
		})
	})
	
	// UI Routes - allow guest access
	uiRoutes := s.router.Group("/")
	uiRoutes.Use(guestMiddleware)
	{
		uiRoutes.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Unixify - Account Management",
			})
		})

		// Register page
		uiRoutes.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"title": "Unixify - Register",
			})
		})
		
		// Claude page - displays documentation content
		uiRoutes.GET("/claude", func(c *gin.Context) {
			// Create documentation content directly
			documentationContent := `# Unixify - UNIX Account/Group Registry

Unixify is a Go application that serves as a registry for UNIX account UIDs and Group GIDs.

## Project Overview

The application provides a web interface for managing UNIX accounts and groups with the following features:

1. PostgreSQL database backend
2. Web interface with four sections: People, System, Database, and Service
3. Complete audit log system for all operations
4. Full RESTful API for all operations
5. JWT-based authentication with optional TOTP 2FA
6. Light/dark mode theme switching with auto-detection
7. Read-only guest mode with visual indicators
8. Gradient text and consistent button styling

## Account/Group Types and UID/GID Ranges

| Type     | Account UID Range | Group GID Range |
|----------|-------------------|-----------------|
| People   | 1000-6000         | 1000-6000       |
| System   | 9000-9100         | 9000-9100       |
| Database | 7000-7999         | 7000-7999       |
| Service  | 8000-8999         | 8000-8999       |

## Key Operations

- Add/edit/delete accounts and groups
- Assign/remove users from groups
- View detailed audit logs of all system events
- Search by UID, GID, username, or groupname
- User authentication with optional TOTP 2FA
- Theme switching (light/dark mode)
- Guest read-only access with registration for edit permissions

## Tech Stack

- Go with Gin web framework
- PostgreSQL database
- HTML/CSS/JavaScript frontend
- RESTful API backend
- JWT-based authentication
- Google Authenticator TOTP support
- Theme switching with CSS variables
- Audit logging for all operations

## Authentication System

The application includes a comprehensive authentication system:
- JWT token-based authentication
- Password hashing with bcrypt
- Optional TOTP second factor with Google Authenticator
- Protected API routes with middleware
- User profiles and password management
- Self-registration with email verification and admin approval
- Automatic guest mode with clear visual indicators
- Proper separation between regular users and guest accounts

## Theming System

The application supports light and dark themes:
- CSS variables for comprehensive theme support
- Theme toggle button integrated in the navigation bar
- Theme preference stored in localStorage for persistence
- System preference detection via prefers-color-scheme
- Dark mode for all UI components including forms, tables, and alerts
- Light grey text in tables for better dark mode readability
- Gradient text effects for headings and descriptions
- Consistent color palette for buttons and interactive elements
- Custom colored badges with theme-appropriate styling

## Access Control

The application implements a role-based access control system:
- Guests (unauthenticated users) have read-only access to view data
- Registration is required to request edit permissions
- New registrations require admin approval
- Authenticated users can perform edits based on their role
- UI dynamically adapts to show/hide edit controls based on permissions
- Clear visual indicators show current access mode (read-only vs. edit mode)`

			c.HTML(http.StatusOK, "document.html", gin.H{
				"title": "Unixify - Documentation",
				"content": documentationContent,
			})
		})

		// Section routes
		sections := []string{"people", "system", "database", "service"}
		for _, section := range sections {
			section := section // Create a new variable to avoid closure issues
			uiRoutes.GET("/"+section, func(c *gin.Context) {
				c.HTML(http.StatusOK, "section.html", gin.H{
					"title":   "Unixify - " + section,
					"section": section,
				})
			})
		}
	}
}

// Run starts the API server
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.config.Server.Port)
	s.logger.Infof("Starting server on %s", addr)
	return s.router.Run(addr)
}