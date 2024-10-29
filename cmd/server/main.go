package main

import (
	"datingApp/db"
	"datingApp/internal/handlers"
	"datingApp/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
    conn, err := db.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    log.Println("Starting migrations...")
    db.RunMigrations(conn)

    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    authHandler := handlers.AuthHandler{DB: conn}
    authHandler.AuthRoutes(r)

    matchHandler := handlers.MatchHandler{
        DB:           conn,
        MatchService: &services.MatchService{DB: conn},
    }
    matchHandler.AuthRoutes(r)

    userHandler := handlers.UserHandler{
        DB:          conn,
        UserService: &services.UserService{DB: conn},
    }
    userHandler.AuthRoutes(r)

    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}