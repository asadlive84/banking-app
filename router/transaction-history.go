package router

import (
	"context"
	"log"
	"net/http"
    "github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

)


func getTransactionHistory(c *gin.Context, mongoCollection *mongo.Collection) {
	accountID := c.Param("account_id")

	cursor, err := mongoCollection.Find(context.TODO(), bson.M{"account_id": accountID})
	if err != nil {
		log.Printf("Failed to fetch transactions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	var transactions []bson.M
	if err := cursor.All(context.TODO(), &transactions); err != nil {
		log.Printf("Failed to parse transactions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
