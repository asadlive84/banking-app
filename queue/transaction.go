package queue

import (
	"context"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"


	"github.com/asadlive84/banking-app/model"


)

func ProcessTransactions(rabbitmqURL string, db *gorm.DB, mongoCollection *mongo.Collection) {
	conn, err := ConnectToRabbitMQ(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"transaction_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		var transaction model.TransactionMessage
		err := json.Unmarshal(msg.Body, &transaction)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		processTransaction(transaction, db, mongoCollection)
	}
}

func processTransaction(transaction model.TransactionMessage, db *gorm.DB, mongoCollection *mongo.Collection) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			tx.Rollback() // Rollback on panic
		}
	}()

	var account model.Account
	if err := tx.Where("account_number = ?", transaction.AccountID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Account not found for account ID: %s", transaction.AccountID)
		} else {
			log.Printf("Failed to fetch account: %v", err)
		}
		tx.Rollback()
		return
	}

	switch transaction.Type {
	case "deposit":
		account.Balance += transaction.Amount
	case "withdraw":
		if account.Balance < transaction.Amount {
			log.Printf("Insufficient balance for withdrawal. Account ID: %s, Balance: %.2f, Amount: %.2f",
				transaction.AccountID, account.Balance, transaction.Amount)
			tx.Rollback()
			return
		}
		account.Balance -= transaction.Amount
	default:
		log.Printf("Invalid transaction type: %s", transaction.Type)
		tx.Rollback()
		return
	}

	if err := tx.Save(&account).Error; err != nil {
		log.Printf("Failed to update account balance in PostgreSQL: %v", err)
		tx.Rollback()
		return
	}

	newTransaction := model.Transaction{
		AccountNumber: transaction.AccountID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		Status:        "completed",
	}
	if err := tx.Create(&newTransaction).Error; err != nil {
		log.Printf("Failed to record transaction in PostgreSQL: %v", err)
		tx.Rollback()
		return
	}

	transactionID, err := primitive.ObjectIDFromHex(transaction.TransactionID)
	if err != nil {
		log.Printf("Invalid transaction ID: %v", err)
		tx.Rollback()
		return
	}

	updateResult, err := mongoCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": transactionID}, // Use the unique transaction ID
		bson.M{"$set": bson.M{"status": "completed"}},
	)
	if err != nil {
		log.Printf("Error updating MongoDB: %v", err)
		tx.Rollback()
		return
	} else {
		log.Printf("MongoDB update result: Matched=%d, Modified=%d", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	if updateResult.MatchedCount == 0 {
		log.Printf("No matching document found in MongoDB for transaction ID: %s", transaction.TransactionID)
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction in PostgreSQL: %v", err)
		tx.Rollback()
		return
	}

	log.Printf("Transaction processed successfully for account ID: %s", transaction.AccountID)
}