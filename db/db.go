package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Пустой импорт для файловых миграций
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB подключается к базе данных и возвращает объект *gorm.DB.
func ConnectDB() (*gorm.DB, error) {
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")

    if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
        return nil, fmt.Errorf("one or more environment variables are missing")
    }

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
        dbHost, dbUser, dbPassword, dbName, dbPort)

    var db *gorm.DB
    var err error
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            log.Println("Connected to the database successfully")
            return db, nil
        }
        log.Printf("Database connection failed: %v. Retrying in 5 seconds...", err)
        time.Sleep(5 * time.Second)
    }
    return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}


// RunMigrations выполняет миграции для базы данных.
func RunMigrations(db *gorm.DB) {
    sqlDB, err := db.DB() // Получаем *sql.DB из *gorm.DB
    if err != nil {
        log.Fatal("Failed to get *sql.DB from *gorm.DB:", err)
    }

    driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{}) // Используем псевдоним
    if err != nil {
        log.Fatal("Failed to create migration driver:", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file:///migrations",
        "datingdb",
        driver,
    )
    if err != nil {
        log.Fatal("Failed to initialize migration:", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("Migrate failed:", err)
    } else {
        log.Println("Migrated successfully")
    }
}