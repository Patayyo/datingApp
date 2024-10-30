package handlers

import (
	"datingApp/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MatchHandler управлеят обработкой запросов, связанных с функцией подбора пользоватлей
type MatchHandler struct {
	DB *gorm.DB // Ссылка на базу данных
	MatchService services.MatchServiceInterface // Интерфейс для работы с сервисом подбора
}

// LikeUser лайкает пользователя
func (h *MatchHandler) LikeUser(c *gin.Context) {
	// Получем идентификатор текущего пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}) // Ошибка, если пользователь не авторизован
		return
	}

	// Извлекаем идентификатор пользователя, которого хотим подобрать
	targetIDStr := c.Param("targetID")
	targetID, err := strconv.Atoi(targetIDStr) // Преобразуем строку в число
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Ошибка, если передан некорректный ID
		return
	}

	// Вызываем метод сервиса подбора, чтобы лайкнуть пользователя
	err = h.MatchService.LikeUser(userID.(int), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like user"})
		return
	}

	// Сообщение об успешном добавлении лайка
	c.JSON(http.StatusOK, gin.H{"message": "User liked successfully"})
}


// GetMatches возвращает список пользователей, которых пользователь лайкнул
func (h *MatchHandler) GetMatches(c *gin.Context) {
	// Получаем идентификатор текущего пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}) // Ошибка, если пользователь не авторизован
		return
	}

	// Вызываем сервис для получения списка мэтчей для пользователя
	matches, err := h.MatchService.GetUserMatches(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get matches"})
		return
	}

	// Отправляем список мэтчей в ответе
	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

// AuthRoutes регистрирует маршруты для обработки запросов, связанных с лайками и матчами.
func (h *MatchHandler) AuthRoutes(router *gin.Engine) {
	// Создаем группу маршрутов с префиксом /match, к которым применяется middleware аутентификации
	match := router.Group("/match")
	match.Use(TokenAuthMiddleware())
	{
		match.POST("/like/:targetID", h.LikeUser) // Маршрут для лайка пользователя
		match.GET("/matches", h.GetMatches) // Маршрут для получения списка матчей
	}
}