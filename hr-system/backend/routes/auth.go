package routes

import (
	"hr-backend/handlers"
	"hr-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		// Public routes
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
		
		// Protected routes
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/change-password", handlers.ChangePassword)
			protected.GET("/profile", handlers.GetProfile)
		}
	}
}