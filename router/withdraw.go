package router

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asadlive84/banking-app/model"
	"github.com/asadlive84/banking-app/queue"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Withdraw(c *gin.Context, transactionCollection *mongo.Collection, rabbitmqURL string) {
	var req struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	transactionID := primitive.NewObjectID()

	transaction := bson.M{
		"_id":        transactionID,
		"account_id": req.AccountID,
		"amount":     req.Amount,
		"type":       "withdraw",
		"status":     "pending",
		"timestamp":  time.Now(),
	}

	_, err := transactionCollection.InsertOne(context.TODO(), transaction)
	if err != nil {
		log.Printf("Failed to insert transaction into MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	queue.PublishToQueue(model.TransactionMessage{
		TransactionID: transactionID.Hex(),
		AccountID:     req.AccountID,
		Amount:        req.Amount,
		Type:          "withdraw",
	}, rabbitmqURL)

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal request submitted"})
}
