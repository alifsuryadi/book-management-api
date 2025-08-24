package routes

import (
	"book-management-api/config"
	"book-management-api/handlers"
	"book-management-api/middleware"
	"book-management-api/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB, cfg *config.Config) {
	// Set trusted proxies (only localhost for development)
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	
	// Add CORS middleware
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "API is healthy",
			Data: gin.H{
				"status":  "ok",
				"service": "book-management-api",
			},
		})
	})

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db, cfg)
	categoryHandler := handlers.NewCategoryHandler(db)
	bookHandler := handlers.NewBookHandler(db)

	// API routes
	api := router.Group("/api")
	{
		// User authentication routes
		users := api.Group("/users")
		{
			users.POST("/login", userHandler.Login)
			users.POST("/seed-admin", userHandler.SeedAdmin) // Temporary endpoint to create admin user
			users.POST("/reset-admin-password", userHandler.ResetAdminPassword) // Temporary endpoint to reset admin password
		}

		// Category routes with JWT authentication
		categories := api.Group("/categories")
		categories.Use(middleware.JWTAuth(cfg)) // Use JWT authentication
		{
			categories.GET("", categoryHandler.GetAll)
			categories.POST("", categoryHandler.Create)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.DELETE("/:id", categoryHandler.Delete)
			categories.GET("/:id/books", categoryHandler.GetBooksByCategory)
		}

		// Book routes with JWT authentication
		books := api.Group("/books")
		books.Use(middleware.JWTAuth(cfg)) // Use JWT authentication
		{
			books.GET("", bookHandler.GetAll)
			books.POST("", bookHandler.Create)
			books.GET("/:id", bookHandler.GetByID)
			books.DELETE("/:id", bookHandler.Delete)
		}
	}

	// Alternative routes with Basic Auth (comment out JWT routes above and uncomment these if you prefer Basic Auth)
	/*
	// Category routes with Basic Authentication
	categories := api.Group("/categories")
	categories.Use(middleware.BasicAuth(cfg))
	{
		categories.GET("", categoryHandler.GetAll)
		categories.POST("", categoryHandler.Create)
		categories.GET("/:id", categoryHandler.GetByID)
		categories.DELETE("/:id", categoryHandler.Delete)
		categories.GET("/:id/books", categoryHandler.GetBooksByCategory)
	}

	// Book routes with Basic Authentication
	books := api.Group("/books")
	books.Use(middleware.BasicAuth(cfg))
	{
		books.GET("", bookHandler.GetAll)
		books.POST("", bookHandler.Create)
		books.GET("/:id", bookHandler.GetByID)
		books.DELETE("/:id", bookHandler.Delete)
	}
	*/

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Route not found",
			Error:   "the requested endpoint does not exist",
		})
	})
}