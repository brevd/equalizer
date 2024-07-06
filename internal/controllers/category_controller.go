package controllers

import (
	"net/http"

	"github.com/brevd/equalizer/internal"
	"github.com/brevd/equalizer/internal/models"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {

	// Find all categories
	rows, err := internal.DB.Query("SELECT id, title, description FROM categories")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Title, &category.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		categories = append(categories, category)
	}

	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := internal.DB.Exec("INSERT INTO categories (title, description) VALUES (?,?)", category.Title, category.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "category created successfully", "category id": categoryID})
}
