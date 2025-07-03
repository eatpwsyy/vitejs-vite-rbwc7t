package routes

import (
	"hr-backend/handlers"
	"hr-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupLeaveRoutes(api *gin.RouterGroup) {
	leaves := api.Group("/leaves")
	leaves.Use(middleware.AuthMiddleware())
	{
		// All authenticated users can view leaves (with filtering)
		leaves.GET("", handlers.GetLeaves)
		leaves.GET("/:id", handlers.GetLeave)
		
		// Employees can create their own leave requests
		leaves.POST("", handlers.CreateLeave)
		
		// Employees can update/delete their own pending leave requests
		leaves.PUT("/:id", handlers.UpdateLeave)
		leaves.DELETE("/:id", handlers.DeleteLeave)
		
		// Managers, HR and Admin can approve/reject leave requests
		leaves.POST("/:id/approve", middleware.RequireRole("admin", "hr", "manager"), handlers.ApproveLeave)
		
		// Get leave types and statuses
		leaves.GET("/types", handlers.GetLeaveTypes)
		leaves.GET("/statuses", handlers.GetLeaveStatuses)
	}
}