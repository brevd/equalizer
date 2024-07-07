package controllers

import (
	"net/http"

	"github.com/brevd/equalizer/internal"
	"github.com/brevd/equalizer/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateBudget(c *gin.Context) {
	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := internal.DB.Exec("INSERT INTO budgets (amount, time_period, category_id, user_id) VALUES (?,?,?,?)", budget.Amount, budget.TimePeriod, budget.CategoryID, budget.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	budgetID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve budget ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Budget created successfully", "budget id": budgetID})
}

func GetBudgetByUserID(c *gin.Context) {
	id := c.Param("id")

	rows, err := internal.DB.Query("SELECT id, category_id, time_period, user_id, amount, created_at, updated_at FROM budgets WHERE user_id = (?)", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var budgets []models.Budget
	for rows.Next() {
		var budget models.Budget
		var createdAt, updatedAt string
		if err := rows.Scan(&budget.ID, &budget.CategoryID, &budget.TimePeriod, &budget.UserID, &budget.Amount, &createdAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		budgets = append(budgets, budget)
	}

	if budgets == nil {
		budgets = []models.Budget{}
	}

	c.JSON(http.StatusOK, budgets)
}
