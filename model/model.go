package model

import "time"

type Account struct {
	ID            int     `gorm:"primaryKey"`
	AccountNumber string  `gorm:"unique;not null"`
	Name          string  `gorm:"not null"`
	Balance       float64 `gorm:"type:numeric(15,2);not null;default:0"`
}

type Transaction struct {
	ID            int       `gorm:"primaryKey"`
	AccountNumber string    `gorm:"foreignKey:AccountNumber;references:AccountNumber;not null"`
	Type          string    `gorm:"check:type IN ('deposit', 'withdraw');not null"`
	Amount        float64   `gorm:"type:numeric(15,2);not null"`
	Timestamp     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status        string    `gorm:"default:'pending'"`
}


type TransactionMessage struct {
	TransactionID string  `json:"transaction_id"`
	AccountID     string  `json:"account_id"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
}