package api

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/home/unixify/internal/config"
	"github.com/home/unixify/internal/handlers"
	"github.com/home/unixify/internal/service"
	"github.com/sirupsen/logrus"
)

// Server represents the API server
type Server struct {
	router  *gin.Engine
	config  *config.Config
	logger  *logrus.Logger
	handler *handlers.Handler
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

	// Create server
	server := &Server{
		router:  router,
		config:  cfg,
		logger:  logger,
		handler: handler,
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
	// API routes
	api := s.router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Account routes
		accounts := api.Group("/accounts")
		{
			accounts.GET("", s.handler.GetAllAccounts)
			accounts.GET("/:id", s.handler.GetAccount)
			accounts.POST("", s.handler.CreateAccount)
			accounts.PUT("/:id", s.handler.UpdateAccount)
			accounts.DELETE("/:id", s.handler.DeleteAccount)
			accounts.GET("/uid/:uid", s.handler.GetAccountByUID)
			accounts.GET("/username/:username", s.handler.GetAccountByUsername)
			accounts.GET("/:id/groups", s.handler.GetAccountGroups)
		}

		// Group routes
		groups := api.Group("/groups")
		{
			groups.GET("", s.handler.GetAllGroups)
			groups.GET("/:id", s.handler.GetGroup)
			groups.POST("", s.handler.CreateGroup)
			groups.PUT("/:id", s.handler.UpdateGroup)
			groups.DELETE("/:id", s.handler.DeleteGroup)
			groups.GET("/gid/:gid", s.handler.GetGroupByGID)
			groups.GET("/groupname/:groupname", s.handler.GetGroupByGroupname)
			groups.GET("/:id/accounts", s.handler.GetGroupMembers)
		}

		// Membership routes
		membership := api.Group("/memberships")
		{
			membership.POST("", s.handler.AssignAccountToGroup)
			membership.DELETE("", s.handler.RemoveAccountFromGroup)
		}

		// Search routes
		search := api.Group("/search")
		{
			search.GET("/accounts", s.handler.SearchAccounts)
			search.GET("/groups", s.handler.SearchGroups)
		}

		// Audit routes
		audit := api.Group("/audit")
		{
			audit.GET("", s.handler.GetAuditEntries)
			audit.GET("/:id", s.handler.GetAuditEntry)
		}
	}

	// Web interface routes
	// Add custom template functions
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	
	// Set HTML templates with functions
	s.router.SetFuncMap(funcMap)
	s.router.LoadHTMLGlob("web/templates/*")
	s.router.Static("/static", "web/static")
	
	// UI Routes
	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Unixify - Account Management",
		})
	})

	// Section routes
	sections := []string{"people", "system", "database", "service"}
	for _, section := range sections {
		s.router.GET("/"+section, func(c *gin.Context) {
			c.HTML(http.StatusOK, "section.html", gin.H{
				"title":   "Unixify - " + section,
				"section": section,
			})
		})
	}
}

// Run starts the API server
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.config.Server.Port)
	s.logger.Infof("Starting server on %s", addr)
	return s.router.Run(addr)
}