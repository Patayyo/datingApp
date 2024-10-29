package handlers

import (
	"datingApp/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatchHandler struct {
	DB *gorm.DB
	MatchService services.MatchServiceInterface
}

func (h *MatchHandler) LikeUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	targetIDStr := c.Param("targetID")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.MatchService.LikeUser(userID.(int), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User liked successfully"})
}

func (h *MatchHandler) GetMatches(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	matches, err := h.MatchService.GetUserMatches(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get matches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

func (h *MatchHandler) AuthRoutes(router *gin.Engine) {
	match := router.Group("/match")
	match.Use(TokenAuthMiddleware())
	{
		match.POST("/like/:targetID", h.LikeUser)
		match.GET("/matches", h.GetMatches)
	}
}