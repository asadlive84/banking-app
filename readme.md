# Banking App

This is a banking app that uses PostgreSQL, MongoDB, and RabbitMQ to handle transaction processing and account management. In this app where data synchronization happens between different databases using messaging queues.

---

## Features

- **Account Creation**: Create new accounts.
- **Transaction Processing**: Handle deposits and withdrawals.
- **Data Synchronization**: Synchronize data between PostgreSQL and MongoDB.
- **Messaging System**: Process transactions using RabbitMQ.

---

## Technology Stack

- **PostgreSQL**: Used for storing account and transaction data.
- **MongoDB**: Used for logging transaction statuses.
- **RabbitMQ**: Used for processing transaction messages.
- **Golang**: The primary programming language.
- **Gin Framework**: Used for building REST APIs.

---

## Installation

### Prerequisites

1. **Go** must be installed. [Download Go](https://golang.org/dl/) I'm using go 1.23 for this project

2. **Docker** must be installed. [Download Docker](https://www.docker.com/)

### Steps

1. Clone the repository:
   ```bash
   git clone git@github.com:asadlive84/banking-app.git
   ```

2. go to project folder:
   ```bash
   cd banking-app
   ```

3. Now you have to run this command for install dependencies for this project:
   ```bash
   go mod tidy
   ```

4. Run this project with Docker compose:
   ```bash
   docker compose up --build 
   ```

5. Now you can look this project structure and you find ```banking-api.postman_collection.json``` this postman collection, and can import this file in your postman client app and test it or you can look this  API Endpoints example

# API Endpoints

## 1. Create a New Account

- **URL**: `/accounts`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "account_number": "1234567890",
    "balance": 1000.00
  }
  ```
- **Response Body**:
```json
  {
  "message": "Account created successfully"
}
```



## 2. Process a Transaction (Deposit)

- **URL**: `/transactions/deposit`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
  "account_id": "1234567890",
  "amount": 500.00,
    }
  ```
- **Response Body**:
```json
  {
  "message": "Deposit request submitted"
}
```


## 3. Process a Transaction (withdraw)

- **URL**: `/transactions/withdraw`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
  "account_id": "1234567890",
  "amount": 200.00,
    }
  ```
- **Response Body**:
```json
  {
  "message": "Withdrawal request submitted"
}
```


## 4. Get account info by account number

- **URL**: `/accounts/ACC12345`
- **Method**: `GET`
- **Response Body**:
```json
  {
    "account_id": "ACC12345",
    "balance": 1499
}
``` 

## 5. Get account transactions info by account number

- **URL**: `/transactions/ACC12345`
- **Method**: `GET`
- **Response Body**:
```json
  [
    {
        "_id": "67c97193dee000d89cee3e24",
        "account_id": "ACC12345",
        "amount": 500,
        "status": "completed",
        "timestamp": "2025-03-06T09:57:39.339Z",
        "type": "deposit"
    },
    {
        "_id": "67c97196dee000d89cee3e25",
        "account_id": "ACC12345",
        "amount": 1,
        "status": "completed",
        "timestamp": "2025-03-06T09:57:42.054Z",
        "type": "withdraw"
    },
    {
        "_id": "67c9738fdee000d89cee3e26",
        "account_id": "ACC12345",
        "amount": 1,
        "status": "completed",
        "timestamp": "2025-03-06T10:06:07.683Z",
        "type": "withdraw"
    },
    {
        "_id": "67c97392dee000d89cee3e27",
        "account_id": "ACC12345",
        "amount": 500,
        "status": "completed",
        "timestamp": "2025-03-06T10:06:10.364Z",
        "type": "deposit"
    }
]
``` 