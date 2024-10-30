package handlers

import (
	"datingApp/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler управляет запросами, связанными с профилем пользователя
type UserHandler struct {
	DB *gorm.DB // Ссылка на базу данных
	UserService *services.UserService // Сервис для управления профилями пользователей
}

// GetProfile возвращает информацию о профиле пользователя
func (h *UserHandler) GetProfile(c * gin.Context) {
	// Извлекаем идентификатор пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthoeized"}) // Ошибка, если пользователль не авторизован
		return
	}

	// Получаем профиль пользователя из сервиса
	profile, err := h.UserService.GetUserProfile(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	// Проверям, найден ли профиль
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"}) // Ошибка, если профиль не найден
		return
	}

	// Возвращаем профиль пользователя в формате JSON
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// UpdateProfile обрабатывает запрос на обновление профиля пользователя
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Извлекаем идентификатор пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorizred"}) // Ошибка, если пользователь не авторизован
		return
	}

	// Получаем данные для обновления профиля из JSON-запроса
	var input services.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем профиль пользователя
	if err := h.UserService.UpdateUserProfile(userID.(int), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Подтверждаем успешное обновление профиля
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

//DeletProfile обрабатывает запрос на удаление профиля пользователя
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	// Получаем идентификатор пользователя из параметров URL
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Удаляем профиль пользователя, используя сервис
	if err := h.UserService.DeleteUserProfile(userID); err != nil {
		// Если произошла ошибка при удалении профиля, возвращаем ошибку
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		}
		return
	}

	// Подтверждаем успешное удаление профиля
	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

// AuthRoutes регистрирует маршруты для запросов, связанных с профилем пользователя.
func (h *UserHandler) AuthRoutes(router *gin.Engine) {
	// Создаем группу маршрутов с префиксом /user, защищенных аутентификацией
	user := router.Group("/user")
	user.Use(TokenAuthMiddleware())
	{
		user.GET("/profile", h.GetProfile) // Маршрут для получения профиля
		user.PUT("/profile", h.UpdateProfile) // Маршрут для обновления профиля
		user.DELETE("/profile/:userID", h.DeleteProfile) // Маршрут для удаления профиля
	}
}