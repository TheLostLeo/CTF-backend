package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thelostleo/CTF-backend/database"
	"github.com/thelostleo/CTF-backend/models"
)

type AdminController struct{}

// CreateChallenge handles POST /admin/challenges
func (ac *AdminController) CreateChallenge(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category" binding:"required"`
		Points      int    `json:"points" binding:"required,min=1"`
		Flag        string `json:"flag" binding:"required"`
		Hint        string `json:"hint"`
		FileURL     string `json:"file_url"`
		IsActive    *bool  `json:"is_active"` // Pointer to handle optional boolean
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Set default value for IsActive if not provided
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	challenge := models.Challenge{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Points:      req.Points,
		Flag:        req.Flag,
		Hint:        req.Hint,
		FileURL:     req.FileURL,
		IsActive:    isActive,
	}

	if err := database.DB.Create(&challenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create challenge",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Challenge created successfully",
		"challenge": gin.H{
			"id":          challenge.ID,
			"title":       challenge.Title,
			"description": challenge.Description,
			"category":    challenge.Category,
			"points":      challenge.Points,
			"hint":        challenge.Hint,
			"file_url":    challenge.FileURL,
			"is_active":   challenge.IsActive,
		},
	})
}

// GetAllUsers handles GET /admin/users
func (ac *AdminController) GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := database.DB.Select("id, username, email, score, is_admin, created_at").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"total_users": len(users),
	})
}

// UpdateChallenge handles PUT /admin/challenges/:id
func (ac *AdminController) UpdateChallenge(c *gin.Context) {
	id := c.Param("id")
	challengeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid challenge ID",
		})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Points      int    `json:"points,omitempty"`
		Flag        string `json:"flag"`
		Hint        string `json:"hint"`
		FileURL     string `json:"file_url"`
		IsActive    *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Check if challenge exists
	var challenge models.Challenge
	if err := database.DB.First(&challenge, challengeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Challenge not found",
		})
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Points > 0 {
		updates["points"] = req.Points
	}
	if req.Flag != "" {
		updates["flag"] = req.Flag
	}
	if req.Hint != "" {
		updates["hint"] = req.Hint
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := database.DB.Model(&challenge).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update challenge",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Challenge updated successfully",
		"challenge": challenge,
	})
}

// DeleteChallenge handles DELETE /admin/challenges/:id
func (ac *AdminController) DeleteChallenge(c *gin.Context) {
	id := c.Param("id")
	challengeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid challenge ID",
		})
		return
	}

	// Soft delete the challenge
	if err := database.DB.Delete(&models.Challenge{}, challengeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete challenge",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Challenge deleted successfully",
	})
}

// GetDashboard handles GET /admin/dashboard
func (ac *AdminController) GetDashboard(c *gin.Context) {
	// Get statistics
	var userCount int64
	var challengeCount int64
	var submissionCount int64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Challenge{}).Where("is_active = ?", true).Count(&challengeCount)
	database.DB.Model(&models.Submission{}).Count(&submissionCount)

	// Get recent submissions
	var recentSubmissions []models.Submission
	database.DB.Preload("User").Preload("Challenge").
		Order("submitted_at DESC").
		Limit(10).
		Find(&recentSubmissions)

	c.JSON(http.StatusOK, gin.H{
		"statistics": gin.H{
			"total_users":       userCount,
			"active_challenges": challengeCount,
			"total_submissions": submissionCount,
		},
		"recent_submissions": recentSubmissions,
	})
}
