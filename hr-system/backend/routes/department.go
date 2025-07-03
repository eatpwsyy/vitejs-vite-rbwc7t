package routes

import (
	"hr-backend/handlers"
	"hr-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDepartmentRoutes(api *gin.RouterGroup) {
	departments := api.Group("/departments")
	departments.Use(middleware.AuthMiddleware())
	{
		// All authenticated users can view departments
		departments.GET("", handlers.GetDepartments)
		departments.GET("/:id", handlers.GetDepartment)
		
		// Only HR and Admin can create, update, delete departments
		departments.POST("", middleware.RequireRole("admin", "hr"), handlers.CreateDepartment)
		departments.PUT("/:id", middleware.RequireRole("admin", "hr"), handlers.UpdateDepartment)
		departments.DELETE("/:id", middleware.RequireRole("admin", "hr"), handlers.DeleteDepartment)
	}
}