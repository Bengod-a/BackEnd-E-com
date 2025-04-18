package db

import (
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/Bengod-a/DB-GO/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	if dsn == "" {
		panic("DATABASE_URL is not set in environment variables")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			// LogLevel:      logger.Info,
			Colorful: true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		PrepareStmt:            false,
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
	})
	if err != nil {
		panic("Failed to connect to Supabase DB: " + err.Error())
	}

	// DB.Migrator().DropTable(
	// 	&models.User{},
	// 	&models.Product{},
	// 	&models.Images{},
	// 	&models.Category1{},
	// 	&models.Category2{},
	// 	&models.Order{},
	// 	&models.Address{},
	// 	&models.Cart{},
	// 	&models.ProductOnCart{},
	// 	&models.Favorite{},
	// 	&models.ProductOnOrder{},
	// )

	// if err := DB.AutoMigrate(
	// 	&models.Product{},
	// 	&models.User{},
	// 	&models.Images{},
	// 	&models.Category1{},
	// 	&models.Category2{},
	// 	&models.Order{},
	// 	&models.Address{},
	// 	&models.Cart{},
	// 	&models.ProductOnCart{},
	// 	&models.Favorite{},
	// 	&models.ProductOnOrder{},
	// ); err != nil {
	// 	panic("Failed to migrate tables: " + err.Error())
	// }

	fmt.Println("Successfully connected DB")
	return DB
}
