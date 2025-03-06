package router

import (
	"net/http"

	"github.com/asadlive84/banking-app/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAccount(c *gin.Context, db *gorm.DB) {
	var req struct {
		Name          string  `json:"name"`
		AccountNumber string  `json:"account_number"`
		Balance       float64 `json:"balance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Balance <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Initial balance cannot be negative or zero"})
		return
	}

	newAccount := model.Account{
		Name:          req.Name,
		AccountNumber: req.AccountNumber,
		Balance:       req.Balance,
	}

	if err := db.Create(&newAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully"})
}
