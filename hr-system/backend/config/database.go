package config

import (
	"log"
	"os"

	"hr-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	
	// Use SQLite for simplicity - can be changed to PostgreSQL/MySQL later
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "hr_system.db"
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate schemas
	err = DB.AutoMigrate(
		&models.Employee{},
		&models.Department{},
		&models.Leave{},
		&models.User{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Seed initial data
	seedData()
}

func seedData() {
	// Create default departments if they don't exist
	departments := []models.Department{
		{Name: "Human Resources", Description: "HR Department"},
		{Name: "Engineering", Description: "Engineering Department"},
		{Name: "Marketing", Description: "Marketing Department"},
		{Name: "Finance", Description: "Finance Department"},
	}

	for _, dept := range departments {
		var existingDept models.Department
		if err := DB.Where("name = ?", dept.Name).First(&existingDept).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				DB.Create(&dept)
				log.Printf("Created department: %s", dept.Name)
			}
		}
	}

	// Create default admin user if it doesn't exist
	var adminUser models.User
	if err := DB.Where("email = ?", "admin@hr.com").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			admin := models.User{
				Email:    "admin@hr.com",
				Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
				Role:     "admin",
				IsActive: true,
			}
			DB.Create(&admin)
			log.Println("Created default admin user: admin@hr.com / password")
		}
	}
}