package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/brevd/equalizer/internal"
	"github.com/brevd/equalizer/internal/auth"
	"github.com/brevd/equalizer/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var login models.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Ensure payment methods is not null
	user.PaymentMethods = []string{"General"}
	user.Email = login.Email

	paymentMethodsJSON, err := json.Marshal(user.PaymentMethods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

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

	result, err := tx.Exec("INSERT INTO users (name, payment_methods, email, password, info) VALUES (?, ?, ?, ?, ?)",
		user.Name, paymentMethodsJSON, user.Email, user.Password, user.Info)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
		return
	}

	// Create a default bill mate
	result, err = tx.Exec("INSERT INTO bill_mates (user_id, name) VALUES (?, ?)", userID, user.Name)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bill mate"})
		return
	}

	billMateID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bill mate ID"})
		return
	}

	// Create a bill group titled "Personal"
	result, err = tx.Exec("INSERT INTO bill_groups (title, description) VALUES (?, ?)", "Personal", "Personal group for user: "+user.Name)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bill group"})
		return
	}

	billGroupID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bill group ID"})
		return
	}

	// Associate the bill mate with the bill group
	_, err = tx.Exec("INSERT INTO bill_mate_to_group (bill_mate_id, bill_group_id) VALUES (?, ?)", billMateID, billGroupID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate bill mate with bill group"})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": userID})
}

func Login(c *gin.Context) {
	var login models.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := internal.DB.QueryRow("SELECT id, email, password FROM users WHERE email = (?)", login.Email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error2": err.Error()})
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	err := auth.AddToBlacklist(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add token to blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func GetUsers(c *gin.Context) {

	// Find all users
	rows, err := internal.DB.Query("SELECT id, name, payment_methods, email, info, created_at, updated_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt, updatedAt, paymentMethodsJSON string
		if err := rows.Scan(&user.ID, &user.Name, &paymentMethodsJSON, &user.Email, &user.Info, &createdAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Parse the Payment Methods string to JSON
		err = json.Unmarshal([]byte(paymentMethodsJSON), &user.PaymentMethods)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment methods"})
			return
		}

		// Parse the time strings into time.Time
		user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse created_at"})
			return
		}
		user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse updated_at"})
			return
		}

		users = append(users, user)
	}

	if users == nil {
		users = []models.User{}
	}

	// Return the list of users
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	row := internal.DB.QueryRow("SELECT * FROM users WHERE id=(?)", id)

	var user models.User
	var createdAt, updatedAt, paymentMethodsJSON string
	if err := row.Scan(&user.ID, &user.Name, &paymentMethodsJSON, &user.Email, &user.Info, &createdAt, &updatedAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the Payment Methods string to JSON
	err := json.Unmarshal([]byte(paymentMethodsJSON), &user.PaymentMethods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment methods"})
		return
	}

	// Parse timestamps
	user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse created_at"})
		return
	}
	user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse updated_at"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Retrieve the existing user from the database
	row := internal.DB.QueryRow("SELECT id, name, payment_methods, email, info, created_at, updated_at FROM users WHERE id=?", id)

	var user models.User
	var createdAt, updatedAt, paymentMethodsJSON string
	if err := row.Scan(&user.ID, &user.Name, &paymentMethodsJSON, &user.Email, &user.Info, &createdAt, &updatedAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the Payment Methods string to JSON
	err := json.Unmarshal([]byte(paymentMethodsJSON), &user.PaymentMethods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment methods"})
		return
	}

	// Parse the time strings into time.Time
	user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse created_at"})
		return
	}
	user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse updated_at"})
		return
	}

	// Bind JSON input to update the user
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user's fields
	if updatedUser.Name != "" {
		user.Name = updatedUser.Name
	}
	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}
	if updatedUser.Info != "" {
		user.Info = updatedUser.Info
	}
	if updatedUser.PaymentMethods != nil {
		user.PaymentMethods = updatedUser.PaymentMethods
	}

	// Marshal the updated payment methods back to JSON
	newPaymentMethodsJSON, err := json.Marshal(user.PaymentMethods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payment methods"})
		return
	}

	// Update the user in the database
	_, err = internal.DB.Exec("UPDATE users SET name=?, payment_methods=?, email=?, info=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		user.Name, newPaymentMethodsJSON, user.Email, user.Info, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	_, err := internal.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "this user no longer exists",
		"userID":  id,
	})
}
