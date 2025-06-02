package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/amito07/ems/internal/config"
	"github.com/amito07/ems/internal/models"
)

var DB *gorm.DB

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")

	// Check if tables exist and only migrate if they don't
	// Since we're using an existing database with predefined schema,
	// we'll skip auto-migration to avoid conflicts
	if !DB.Migrator().HasTable(&models.Student{}) {
		log.Println("Running auto migrations...")
		err = AutoMigrate()
		if err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		log.Println("Database migrations completed successfully")
	} else {
		log.Println("Using existing database schema")
	}

	return nil
}

// AutoMigrate runs auto migration for all models
func AutoMigrate() error {
	// Since we're using existing database schema, we'll just verify the connection
	// and check if tables exist. GORM will work with existing tables.
	
	// Check if the main tables exist
	if !DB.Migrator().HasTable(&models.Student{}) {
		log.Println("Warning: students table does not exist, creating...")
		if err := DB.AutoMigrate(&models.Student{}); err != nil {
			return err
		}
	}
	
	if !DB.Migrator().HasTable(&models.Teacher{}) {
		log.Println("Warning: teachers table does not exist, creating...")
		if err := DB.AutoMigrate(&models.Teacher{}); err != nil {
			return err
		}
	}
	
	if !DB.Migrator().HasTable(&models.Course{}) {
		log.Println("Warning: courses table does not exist, creating...")
		if err := DB.AutoMigrate(&models.Course{}); err != nil {
			return err
		}
	}
	
	if !DB.Migrator().HasTable(&models.Enrollment{}) {
		log.Println("Warning: enrollments table does not exist, creating...")
		if err := DB.AutoMigrate(&models.Enrollment{}); err != nil {
			return err
		}
	}
	
	log.Println("Using existing database schema")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
