package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/asadlive84/banking-app/config"
	database "github.com/asadlive84/banking-app/db"
	"github.com/asadlive84/banking-app/queue"
	"github.com/asadlive84/banking-app/router"
)

func main() {
	c, confErr := config.LoadConfig("config")
	if confErr != nil {
		log.Fatalf("Failed to load configuration: %v", confErr)
		return
	}

	db, err := database.InitDB(c.DATA_SOURCE_URL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return
	}

	transactionCollection, mErr := database.InitMongoDB(c.MONGO_SOURCE_URL)
	if mErr != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", mErr)
		return
	}

	conn, err := queue.ConnectToRabbitMQ(c.RABBIT_SOURCE_URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	go queue.ProcessTransactions(c.RABBIT_SOURCE_URL, db, transactionCollection)

	r := router.SetupRouter(db, transactionCollection, c.RABBIT_SOURCE_URL)

	log.Printf("Starting server on port %s...", c.APPLICATION_PORT)
	r.Run(fmt.Sprintf(":%s", c.APPLICATION_PORT))
}
