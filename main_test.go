package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asadlive84/banking-app/model"
	"github.com/asadlive84/banking-app/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRouter *gin.Engine

func TestMain(m *testing.M) {
	var err error
	testDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err := testDB.AutoMigrate(&model.Account{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	testRouter = setupTestRouter(testDB)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.POST("/accounts", func(c *gin.Context) {
		router.CreateAccount(c, db)
	})
	r.GET("/accounts/:id", func(c *gin.Context) {
		router.GetAccountBalance(c, db)
	})
	return r
}

func TestCreateAccount(t *testing.T) {
	reqBody := map[string]interface{}{
		"name":           "John Doe",
		"account_number": "123456",
		"balance":        1000.0,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Account created successfully", response["message"])

	var account model.Account
	err = testDB.Where("account_number = ?", "123456").First(&account).Error
	assert.Nil(t, err)
	assert.Equal(t, "John Doe", account.Name)
	assert.Equal(t, 1000.0, account.Balance)
}

type AccountResponse struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

func TestGetAccount(t *testing.T) {
	account := model.Account{Name: "Jane Doe", AccountNumber: "654321", Balance: 500.0}
	testDB.Create(&account)

	req, _ := http.NewRequest("GET", "/accounts/654321", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := AccountResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "654321", response.AccountID)
	assert.Equal(t, 500.0, response.Balance)
}
