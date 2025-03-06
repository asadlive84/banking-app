package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, mongoCollection *mongo.Collection, rabbitmqURL string) *gin.Engine {
	r := gin.Default()
	r.POST("/accounts", func(c *gin.Context) {
		CreateAccount(c, db)
	})
	r.POST("/transactions/deposit", func(c *gin.Context) {
		Deposit(c, mongoCollection, rabbitmqURL)
	})
	r.POST("/transactions/withdraw", func(c *gin.Context) {
		Withdraw(c, mongoCollection, rabbitmqURL)
	})
	r.GET("/accounts/:id", func(c *gin.Context) {
		GetAccountBalance(c, db)
	})
	r.GET("/transactions/:account_id", func(c *gin.Context) {
		getTransactionHistory(c, mongoCollection)
	})
	return r
}
