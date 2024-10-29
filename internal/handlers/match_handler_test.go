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

type MockMatchService struct {
	mock.Mock
}

func (m *MockMatchService) LikeUser(userID, targetID int) error {
	args := m.Called(userID, targetID)
	return args.Error(0)
}

func (m *MockMatchService) GetUserMatches(userID int) ([]model.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.User), args.Error(1)
}

var _ services.MatchServiceInterface = (*MockMatchService)(nil)

func TestLikeUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockService := new(MockMatchService)
	MatchHandler := &MatchHandler{
		MatchService: mockService,
	}

	router.POST("/match/like/:targetID", func(c *gin.Context) {
		c.Set("userID", 1)
		MatchHandler.LikeUser(c)
	})

	targetID := 2
	mockService.On("LikeUser", 1, targetID).Return(nil)
	req, _ := http.NewRequest("POST", "/match/like/"+strconv.Itoa(targetID), bytes.NewBuffer(nil))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"message": "User liked successfully"}`, recorder.Body.String())

	mockService.AssertCalled(t, "LikeUser", 1, targetID)
}
