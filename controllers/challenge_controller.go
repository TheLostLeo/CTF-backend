package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thelostleo/CTF-backend/database"
	"github.com/thelostleo/CTF-backend/models"
)

type ChallengeController struct{}

// GetAllChallenges handles GET /challenges
func (cc *ChallengeController) GetAllChallenges(c *gin.Context) {
	var challenges []models.Challenge

	// Only show active challenges and hide the flag
	if err := database.DB.Select("id, title, description, category, points, hint, is_active, file_url, created_at").
		Where("is_active = ?", true).
		Find(&challenges).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch challenges",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"challenges":       challenges,
		"total_challenges": len(challenges),
	})
}

// GetChallengeByID handles GET /challenges/:id
func (cc *ChallengeController) GetChallengeByID(c *gin.Context) {
	id := c.Param("id")
	challengeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid challenge ID",
		})
		return
	}

	var challenge models.Challenge
	if err := database.DB.Select("id, title, description, category, points, hint, is_active, file_url, created_at").
		Where("id = ? AND is_active = ?", challengeID, true).
		First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Challenge not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"challenge": challenge,
	})
}

// SubmitFlag handles POST /challenges/:id/submit
func (cc *ChallengeController) SubmitFlag(c *gin.Context) {
	id := c.Param("id")
	challengeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid challenge ID",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Parse request body
	var req struct {
		Flag string `json:"flag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Flag is required",
		})
		return
	}

	// Get challenge details
	var challenge models.Challenge
	if err := database.DB.Where("id = ? AND is_active = ?", challengeID, true).
		First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Challenge not found",
		})
		return
	}

	// Check if user already solved this challenge
	var existingSubmission models.Submission
	if err := database.DB.Where("user_id = ? AND challenge_id = ? AND is_correct = ?",
		userID, challengeID, true).First(&existingSubmission).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Challenge already solved",
		})
		return
	}

	// Check if flag is correct
	isCorrect := req.Flag == challenge.Flag

	// Create submission record
	submission := models.Submission{
		UserID:      userID.(uint),
		ChallengeID: uint(challengeID),
		Flag:        req.Flag,
		IsCorrect:   isCorrect,
		IPAddress:   c.ClientIP(),
		SubmittedAt: time.Now(),
	}

	if err := database.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to record submission",
		})
		return
	}

	// If correct, update user score
	if isCorrect {
		if err := database.DB.Model(&models.User{}).
			Where("id = ?", userID).
			Update("score", database.DB.Raw("score + ?", challenge.Points)).Error; err != nil {
			// Log error but don't fail the request
			c.JSON(http.StatusOK, gin.H{
				"correct": true,
				"message": "Correct flag! Points awarded.",
				"points":  challenge.Points,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"correct": true,
			"message": "Correct flag! Points awarded.",
			"points":  challenge.Points,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"correct": false,
			"message": "Incorrect flag. Try again!",
		})
	}
}
