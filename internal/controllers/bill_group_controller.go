package controllers

import (
	"net/http"

	"github.com/brevd/equalizer/internal"
	"github.com/brevd/equalizer/internal/models"
	"github.com/gin-gonic/gin"
)

func GetBillGroups(c *gin.Context) {

	// Find all users
	rows, err := internal.DB.Query("SELECT id, title, description FROM bill_groups")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var billGroups []models.BillGroup
	for rows.Next() {
		var billGroup models.BillGroup
		if err := rows.Scan(&billGroup.ID, &billGroup.Title, &billGroup.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		billGroups = append(billGroups, billGroup)
	}

	if billGroups == nil {
		billGroups = []models.BillGroup{}
	}

	// Return the list of users
	c.JSON(http.StatusOK, billGroups)
}
