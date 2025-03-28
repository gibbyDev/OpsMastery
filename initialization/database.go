package initialization

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/gibbyDev/OpsMastery/models"
)

func SetupDatabase() *gorm.DB {

    dbHost := GetEnv("DB_HOST")
    dbUser := GetEnv("DB_USER")
    dbPassword := GetEnv("DB_PASSWORD")
    dbName := GetEnv("DB_NAME")
    dbPort := GetEnv("DB_PORT")
    dbSSLMode := GetEnv("DB_SSLMODE")
    dbTimeZone := GetEnv("DB_TIMEZONE")

    dsn := "host=" + dbHost +
        " user=" + dbUser +
        " password=" + dbPassword +
        " dbname=" + dbName +
        " port=" + dbPort +
        " sslmode=" + dbSSLMode +
        " TimeZone=" + dbTimeZone

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    err = db.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatalf("failed to auto migrate User model: %v", err)
    }

    err = db.AutoMigrate(&models.Ticket{})
    if err != nil {
        log.Fatalf("failed to auto migrate Ticket model: %v", err)
    }

    err = db.AutoMigrate(&models.Client{})
    if err != nil {
        log.Fatalf("failed to auto migrate Client model: %v", err)
    }

    return db
}