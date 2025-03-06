package router

import (
	"net/http"

	"github.com/asadlive84/banking-app/model"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func GetAccountBalance(c *gin.Context, db *gorm.DB) {
	accountID := c.Param("id")

	var account model.Account
	if err := db.Where("account_number = ?", accountID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"balance":    account.Balance,
	})
}
