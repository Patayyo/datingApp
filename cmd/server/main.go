package main

import (
	"datingApp/db"                // Пакет для работы с базой данных
	"datingApp/internal/handlers" // Пакет с обработчиком запросов
	"datingApp/internal/services" // Пакет с логикой бизнес-уровня(сервисами)
	"log"

	"github.com/gin-gonic/gin" // Веб-фреймворк Gin для создания HTTP API
)

func main() {
    // Подключение к базе данных
    conn, err := db.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err) // Логирование и завершение программы при ошибке
    }
    
    log.Println("Starting migrations...")
    db.RunMigrations(conn) // Запуск миграций для настройки таблиц в базе данных

    // Создание экземпляра маршрутизатора Gin
    r := gin.Default()

    // Маршрут для проверки доступности сервера
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong", // Возвращает JSON ответ с сообщением pong
        })
    })

    // Инициализация обработчика для маршрутов аутентификации
    authHandler := handlers.AuthHandler{DB: conn}
    authHandler.AuthRoutes(r) // Добавление маршрутов аутентификации в маршрутизатор

    // Инициализация обработчика для маршрутов связей(мэтчей) между пользователями
    matchHandler := handlers.MatchHandler{
        DB:           conn, // Передача подключения к БД
        MatchService: &services.MatchService{DB: conn}, // Передача экземпляра сервиса MatchService
    }
    matchHandler.AuthRoutes(r) // Добавление маршрутов для работы с мэтчами в маршрутизатор

    // Инициализация обработчика для маршрутов работы с пользователями
    userHandler := handlers.UserHandler{
        DB:          conn,
        UserService: &services.UserService{DB: conn},
    }
    userHandler.AuthRoutes(r)

    // Запуск HTTP-сервера на порту 8080 и логирование ошибок, если запуск не удался
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}