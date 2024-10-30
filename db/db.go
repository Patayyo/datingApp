package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"                                   // Библиотека для управления миграциями базы данных
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres" // Драйвер миграций для PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"                     // Подключение для миграций из файловой системы
	"gorm.io/driver/postgres"                                                // Драйвер GORM для PostgreSQL
	"gorm.io/gorm"                                                           // ORM-библиотека GORM
)

// ConnectDB подключается к базе данных и возвращает объект *gorm.DB.
func ConnectDB() (*gorm.DB, error) {
    // Чтение переменных окружения для конфигурации подключения к базе данных.
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")

    // Проверка наличия всех переменных окружения, необходимых для подключения.
    if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
        return nil, fmt.Errorf("one or more environment variables are missing")
    }

    // Формирование строки подключения DSN для PostgreSQL
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
        dbHost, dbUser, dbPassword, dbName, dbPort)

    var db *gorm.DB
    var err error
    maxRetries := 5 // Максимальное количество попыток подключения
    for i := 0; i < maxRetries; i++ {
        // Попытка подключения к базе данных
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            log.Println("Connected to the database successfully")
            return db, nil // Успешное подключение возвращает объект *gorm.DB
        }
        log.Printf("Database connection failed: %v. Retrying in 5 seconds...", err)
        time.Sleep(5 * time.Second) // Задержка перед следующей попыткой подключения
    }
    return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}


// RunMigrations выполняет миграции для базы данных.
func RunMigrations(db *gorm.DB) {
    sqlDB, err := db.DB() // Получаем *sql.DB из *gorm.DB
    if err != nil {
        log.Fatal("Failed to get *sql.DB from *gorm.DB:", err)
    }

    // Создаем драйвер миграции для базы данных
    driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{}) // Используем псевдоним
    if err != nil {
        log.Fatal("Failed to create migration driver:", err)
    }

    // Создаем объект для миграций, указывая путь к миграциям и базу данных
    m, err := migrate.NewWithDatabaseInstance(
        "file:///migrations", // Путь к миграциям
        "datingdb", //  Имя базы данных
        driver,
    )
    if err != nil {
        log.Fatal("Failed to initialize migration:", err)
    }

    // Запуск миграций. Если миграции выполнены, и нет изменейни, просто завершаем. 
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("Migrate failed:", err)
    } else {
        log.Println("Migrated successfully")
    }
}