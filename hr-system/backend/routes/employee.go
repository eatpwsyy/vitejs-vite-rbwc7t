package routes

import (
	"hr-backend/handlers"
	"hr-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEmployeeRoutes(api *gin.RouterGroup) {
	employees := api.Group("/employees")
	employees.Use(middleware.AuthMiddleware())
	{
		// All authenticated users can view employees
		employees.GET("", handlers.GetEmployees)
		employees.GET("/:id", handlers.GetEmployee)
		
		// Only HR and Admin can create, update, delete employees
		employees.POST("", middleware.RequireRole("admin", "hr"), handlers.CreateEmployee)
		employees.PUT("/:id", middleware.RequireRole("admin", "hr"), handlers.UpdateEmployee)
		employees.DELETE("/:id", middleware.RequireRole("admin", "hr"), handlers.DeleteEmployee)
	}
}