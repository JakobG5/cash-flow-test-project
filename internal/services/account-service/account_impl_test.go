package accountservice

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/models"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *sql.DB
var testQueries *db.Queries
var testService IAccountService

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("postgres", "postgres://cashflow_user:cashflow_pass@localhost:5432/cashflow_dev?sslmode=disable")
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	if err := testDB.Ping(); err != nil {
		panic("Failed to ping test database: " + err.Error())
	}

	testQueries = db.New(testDB)
	testLogger := loggermanager.NewLogger("debug")
	testConfig := &models.Config{APIKeyHash: "test_hash_key_123456789"}
	testService = NewAccountService(testQueries, testLogger, testConfig)

	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func TestCreateMerchant_Success(t *testing.T) {
	name := "Test Merchant"
	email := "test@example.com"

	response, err := testService.CreateMerchant(name, email)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Status)
	assert.Equal(t, name, response.Name)
	assert.Equal(t, email, response.Email)
	assert.NotEmpty(t, response.MerchantID)
	assert.NotEmpty(t, response.APIKey)
	assert.Contains(t, response.MerchantID, "CASM-")
	assert.Contains(t, response.APIKey, "api_")
	assert.Equal(t, "Merchant created successfully", response.Message)

	merchant, err := testQueries.GetMerchantWithAPIKey(context.Background(), response.MerchantID)
	require.NoError(t, err)
	assert.Equal(t, name, merchant.Name)
	assert.Equal(t, email, merchant.Email)
	assert.Equal(t, response.MerchantID, merchant.MerchantID)

	_, err = testDB.Exec("DELETE FROM merchants WHERE merchant_id = $1", response.MerchantID)
	require.NoError(t, err)
}

func TestCreateMerchant_DuplicateEmail(t *testing.T) {
	name1 := "First Merchant"
	email := "duplicate@example.com"

	response1, err := testService.CreateMerchant(name1, email)
	require.NoError(t, err)
	require.NotNil(t, response1)

	name2 := "Second Merchant"
	response2, err := testService.CreateMerchant(name2, email)

	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateEmail, err)
	assert.Nil(t, response2)

	_, err = testDB.Exec("DELETE FROM merchants WHERE merchant_id = $1", response1.MerchantID)
	require.NoError(t, err)
}

func TestGetMerchantByID_Success(t *testing.T) {
	name := "Retrieve Test Merchant"
	email := "retrieve@example.com"

	createResponse, err := testService.CreateMerchant(name, email)
	require.NoError(t, err)
	require.NotNil(t, createResponse)

	getResponse, err := testService.GetMerchantByID(createResponse.MerchantID)

	require.NoError(t, err)
	assert.NotNil(t, getResponse)
	assert.True(t, getResponse.Status)
	assert.Equal(t, createResponse.MerchantID, getResponse.MerchantID)
	assert.Equal(t, name, getResponse.Name)
	assert.Equal(t, email, getResponse.Email)
	assert.NotEmpty(t, getResponse.APIKey)
	assert.Equal(t, "Merchant details retrieved successfully", getResponse.Message)

	_, err = testDB.Exec("DELETE FROM merchants WHERE merchant_id = $1", createResponse.MerchantID)
	require.NoError(t, err)
}

func TestGetMerchantByID_NotFound(t *testing.T) {
	nonExistentID := "CASM-NONEXISTENT123"

	response, err := testService.GetMerchantByID(nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "merchant not found")
}

func TestCreateMerchant_NameWithNumbers(t *testing.T) {
	name := "Test123"
	email := "numbers@example.com"

	response, err := testService.CreateMerchant(name, email)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Status)
	assert.Equal(t, name, response.Name)
	assert.Equal(t, email, response.Email)

	_, err = testDB.Exec("DELETE FROM merchants WHERE merchant_id = $1", response.MerchantID)
	require.NoError(t, err)
}

func TestCreateMerchant_EmptyName(t *testing.T) {
	name := ""
	email := "emptyname@example.com"

	response, err := testService.CreateMerchant(name, email)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Status)
	assert.Equal(t, name, response.Name)
	assert.Equal(t, email, response.Email)

	_, err = testDB.Exec("DELETE FROM merchants WHERE merchant_id = $1", response.MerchantID)
	require.NoError(t, err)
}

func TestCreateMerchant_InvalidEmailFormat(t *testing.T) {
	name := "Valid Name"
	email := "invalid-email-format-123"

	response, err := testService.CreateMerchant(name, email)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Status)
	assert.Equal(t, name, response.Name)
	assert.Equal(t, email, response.Email)

	_, err = testDB.Exec("DELETE FROM merchants WHERE merchant_id = $1", response.MerchantID)
	require.NoError(t, err)
}
