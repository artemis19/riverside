package api

import (
	"fmt"
	"net/http"

	"crypto/sha256"
	"encoding/hex"
	"github.com/artemis19/viz/server/database"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var json RegisterRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userAlreadyExists database.User

	// Hash password server-side
	hasher := sha256.New()
	hasher.Write([]byte(json.Password))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))

	query := database.DB.Where("username = ? AND password_hash = ?", json.Username, passwordHash).First(&userAlreadyExists)
	if query.Error != nil {
		// User doesn't exist, create it
		user := &database.User{
			Username:     json.Username,
			PasswordHash: passwordHash,
		}
		database.DB.Create(user)

		c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("Successfully registered with username '%s'!", json.Username)})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Username '%s' already exists!", json.Username)})
	}
}
