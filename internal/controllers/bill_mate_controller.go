package controllers

import (
	"net/http"

	"github.com/brevd/equalizer/internal"
	"github.com/brevd/equalizer/internal/models"
	"github.com/gin-gonic/gin"
)

func GetBillMates(c *gin.Context) {

	// Find all users
	rows, err := internal.DB.Query("SELECT id, user_id, name FROM bill_mates")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var billMates []models.BillMate
	for rows.Next() {
		var billMate models.BillMate
		if err := rows.Scan(&billMate.ID, &billMate.UserID, &billMate.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		billMates = append(billMates, billMate)
	}

	if billMates == nil {
		billMates = []models.BillMate{}
	}

	// Return the list of users
	c.JSON(http.StatusOK, billMates)
}

func GetBillMateById(c *gin.Context) {
	id := c.Param("id")

	row := internal.DB.QueryRow("SELECT id, user_id, name FROM bill_mates WHERE id = (?)", id)

	var billMate models.BillMate
	if err := row.Scan(&billMate.ID, &billMate.UserID, &billMate.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billMate)
}

func CreateBillMate(c *gin.Context) {
	var billMate models.BillMate
	if err := c.ShouldBindJSON(&billMate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := internal.DB.Exec("INSERT INTO bill_mates (user_id, name) VALUES (?,?)", nil, billMate.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	billMateID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "bill mate created successfully", "bill mate": billMateID})
}
