package initialization

import (
	"log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/gibbyDev/OpsMastery/models"
)

func SetupDatabase() *gorm.DB {

	dsn := "host=" + GetEnv("DB_HOST") +
		" user=" + GetEnv("DB_USER") +
		" password=" + GetEnv("DB_PASSWORD") +
		" dbname=" + GetEnv("DB_NAME") +
		" port=" + GetEnv("DB_PORT") +
		" sslmode=" + GetEnv("DB_SSLMODE") +
		" TimeZone=" + GetEnv("DB_TIMEZONE")

	log.Println("Connecting to database with DSN:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	err = db.AutoMigrate(&models.Ticket{})
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	err = db.AutoMigrate(&models.Client{})
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	log.Println("Database connected successfully!")
	return db
} 