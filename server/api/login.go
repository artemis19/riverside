package api

import (
	"fmt"
	"net/http"

	"crypto/sha256"
	"encoding/hex"
	"github.com/artemis19/viz/server/database"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var json LoginRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExists database.User

	// Hash password server-side
	hasher := sha256.New()
	hasher.Write([]byte(json.Password))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))

	query := database.DB.Where("username = ? AND password_hash = ?", json.Username, passwordHash).First(&userExists)
	if query.Error != nil {
		// User doesn't exist or incorrect credentials
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Username or password is incorrect.")})

	} else {
		// User exists, log them in
		database.DB.Save(&userExists)

		c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("Successfully logged in as user '%s'!", json.Username)})
	}
}
