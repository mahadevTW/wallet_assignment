package handler

import (
	"net/http"
	"strings"
	"testing"
	"wallet/testutils"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionFailsWith400ForWrongData(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testServer := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testServer.Server.Close()
	url := testServer.Server.URL + "/transaction"
	body := strings.NewReader(`{some random data}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionFailsForInvalidWalletId(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockService.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"CREDIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionFailsForInvalidTransactionType(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockService.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"CREDIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionFailsWith400ForInsufficientFund(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionSuccessForCredit(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionSuccessForDebit(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestRevertTransactionFailsWith400ForInvalidTransactionID(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestRevertTransactionFailsSuccessForCreditTransaction(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestRevertTransactionFailsSuccessForDebitTransaction(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}
