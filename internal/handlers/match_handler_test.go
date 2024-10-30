package handlers

import (
	"bytes"
	"datingApp/internal/services"
	"datingApp/pkg/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMatchService - структура для создания моков сервиса совпадений.
type MockMatchService struct {
	mock.Mock // Встраивание функциональности mock
}

// LikeUser - метод, который используется для имитации функции "нравится" у сервиса совпадений.
func (m *MockMatchService) LikeUser(userID, targetID int) error {
	args := m.Called(userID, targetID) // Получаем аргументы, переданные в метод
	return args.Error(0) // Возвращаем ошибку, если она есть 
}

// GetUserMatches - метод, который используется для имитации функции получения совпадений пользователя.
func (m *MockMatchService) GetUserMatches(userID int) ([]model.User, error) {
	args := m.Called(userID) // Получаем аргументы, переданные в метод
	return args.Get(0).([]model.User), args.Error(1) // Возвращаем список пользователей и ошибку
}

// Убедимся, что MockMatchService реализует интерфейс MatchServiceInterface
var _ services.MatchServiceInterface = (*MockMatchService)(nil)

// TestLikeUser - тест для проверки функциональности "нравится" пользователя.
func TestLikeUser(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем режим тестирования Gin
	router := gin.Default() // Создаем новый роутер

	mockService := new(MockMatchService) // Создаем экземпляр MockMatchService
	MatchHandler := &MatchHandler{ // Инициализируем MatchHandler с мок-сервисом
		MatchService: mockService,
	}

	// Определяем маршрут для лайка пользователя 
	router.POST("/match/like/:targetID", func(c *gin.Context) {
		c.Set("userID", 1) // Устанавливаем ID пользователя в контекст
		MatchHandler.LikeUser(c) // Вызываем метод лайка пользователя
	})

	targetID := 2 // ID целевого пользователя, которому ставится лайк
	mockService.On("LikeUser", 1, targetID).Return(nil) // Настраиваем ожидание вызова метода
	req, _ := http.NewRequest("POST", "/match/like/"+strconv.Itoa(targetID), bytes.NewBuffer(nil)) // Создаем новый HTTP запрос
	recorder := httptest.NewRecorder() // Создаем новый рекордер для записи ответов
	router.ServeHTTP(recorder, req) // Выполняем HTTP запрос через роутер

	// Проверяем, что код ответа соответствует ожидаемому
	assert.Equal(t, http.StatusOK, recorder.Code)
	// Проверяем, что ответ в формате JSON соответствует ожидаемому
	assert.JSONEq(t, `{"message": "User liked successfully"}`, recorder.Body.String())

	// Проверяем, что метод LikeUser был вызван с ожидаемыми аргументами
	mockService.AssertCalled(t, "LikeUser", 1, targetID)
}
