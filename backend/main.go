package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	initDB()

	// Setup router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	setupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("attendance.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schemas
	err = db.AutoMigrate(&User{}, &Attendance{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized successfully")
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", login)
			auth.POST("/register", register)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware())
		{
			protected.GET("/profile", getProfile)
			protected.PUT("/profile", updateProfile)
			
			// Attendance routes
			protected.POST("/attendance/checkin", checkIn)
			protected.POST("/attendance/checkout", checkOut)
			protected.GET("/attendance", getAttendanceHistory)
			protected.GET("/attendance/today", getTodayAttendance)
			protected.GET("/attendance/stats", getAttendanceStats)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(authMiddleware(), adminMiddleware())
		{
			admin.GET("/users", getAllUsers)
			admin.GET("/attendance", getAllAttendance)
			admin.PUT("/users/:id", updateUser)
			admin.DELETE("/users/:id", deleteUser)
		}
	}
}