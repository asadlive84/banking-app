package db

import (
	"context"
	"log"

	"github.com/asadlive84/banking-app/model"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dataSource string) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(&model.Account{}, &model.Transaction{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return nil, err
	}

	return db, nil
}

func InitMongoDB(dataSource string) (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(dataSource)
	var err error
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	transactionCollection := mongoClient.Database("banking_ledger").Collection("transactions")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"account_id": 1},
		Options: options.Index().SetUnique(false),
	}
	_, err = transactionCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Printf("Failed to create index on account_id: %v", err)
		return nil, err
	}
	return transactionCollection, nil
}
