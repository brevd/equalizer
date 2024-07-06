package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brevd/equalizer/internal"
	"github.com/brevd/equalizer/internal/models"
	"github.com/gin-gonic/gin"
)

func GetExpenses(c *gin.Context) {
	// Find all expenses
	rows, err := internal.DB.Query("SELECT title, amount, description, date, payment_method, vendor, user_id, bill_group_id, category_id FROM expenses")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		var date string
		if err := rows.Scan(&expense.Title, &expense.Amount, &expense.Description, &date, &expense.PaymentMethod, &expense.Vendor, &expense.UserID, &expense.BillGroupID, &expense.CategoryID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Parse the time strings into time.Time
		expense.Date, err = time.Parse(time.RFC3339, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse created_at"})
			return
		}

		expenses = append(expenses, expense)
	}

	if expenses == nil {
		expenses = []models.Expense{}
	}

	c.JSON(http.StatusOK, expenses)
}

func GetExpenseByID(c *gin.Context) {
	id := c.Param("id")

	row := internal.DB.QueryRow("SELECT id, title, amount, description, date, payment_method, vendor, user_id, bill_group_id, category_id FROM expenses WHERE id = (?)", id)

	rows, err := internal.DB.Query("SELECT id, paid, responsible, bill_mate_id, expense_id FROM splits WHERE expense_id = (?)", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var expense models.Expense
	var date string
	err = row.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Description, &date, &expense.PaymentMethod, &expense.Vendor, &expense.UserID, &expense.BillGroupID, &expense.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the time strings into time.Time
	expense.Date, err = time.Parse(time.RFC3339, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse created_at"})
		return
	}

	var splits []models.Split
	for rows.Next() {
		var split models.Split
		if err := rows.Scan(&split.ID, &split.Paid, &split.Responsible, &split.BillMateID, &split.ExpenseID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		splits = append(splits, split)
	}

	if splits == nil {
		splits = []models.Split{}
	}

	completeExpense := models.ExpenseWithSplits{
		Expense: expense,
		Splits:  splits,
	}

	c.JSON(http.StatusOK, completeExpense)
}

func CreateExpense(c *gin.Context) {
	var completeExpense models.ExpenseWithSplits
	if err := c.ShouldBindJSON(&completeExpense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()

	// Validate total paid and responsible amounts
	var totalPaid, totalResponsible int
	for _, split := range completeExpense.Splits {
		totalPaid += split.Paid
		totalResponsible += split.Responsible
	}

	if totalPaid != completeExpense.Expense.Amount || totalResponsible != completeExpense.Expense.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Total paid (%d) and responsible amounts (%d) must equal the expense amount (%d)", totalPaid, totalResponsible, completeExpense.Expense.Amount),
		})
		return
	}

	tx, err := internal.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		}
	}()

	// Insert the expense
	result, err := tx.Exec("INSERT INTO expenses (title, amount, description, date, payment_method, vendor, bill_group_id, user_id, category_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		completeExpense.Expense.Title, completeExpense.Expense.Amount, completeExpense.Expense.Description, now, completeExpense.Expense.PaymentMethod, completeExpense.Expense.Vendor, completeExpense.Expense.BillGroupID, completeExpense.Expense.UserID, completeExpense.Expense.CategoryID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expenseID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expense ID"})
		return
	}

	stmt, err := tx.Prepare("INSERT INTO splits (expense_id, bill_mate_id, paid, responsible) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer stmt.Close()

	// Insert the splits using the prepared statement
	for _, split := range completeExpense.Splits {
		_, err = stmt.Exec(expenseID, split.BillMateID, split.Paid, split.Responsible)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create split"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully", "expense_id": expenseID})
}
