package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"io/ioutil"

	"fmt"
	"net/http"
	"testing"
	"wallet/app/model"
	"wallet/testutils"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
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

func TestCreateTransactionErrorForInvalid(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"Invalid"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionFailsForInvalidTransactionType(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockService.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DUMMY"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionFailsWith400ForInsufficientFund(t *testing.T) {
	mockDatabase := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockDatabase.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	columns := []string{"ID", "Balance"}
	mockDatabase.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 200.0))
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateTransactionSuccessForCredit(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockService.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"CREDIT"}`)
	columns := []string{"ID", "balance"}
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(123, 1000.0))
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(123, 1000.0))
	mockService.Mock.ExpectBegin()
	mockService.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockService.Mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	mockService.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockService.Mock.ExpectCommit()
}
func TestCreateTransactionFailsForWithQueryFailsTOUpdateWallet(t *testing.T) {
	mockDatabase := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockDatabase.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	columns := []string{"ID", "Balance"}
	mockDatabase.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 200.0))
	mockDatabase.Mock.ExpectBegin()
	mockDatabase.Mock.ExpectQuery("UPDATE").WillReturnError(fmt.Errorf("Something went wrong"))
	mockDatabase.Mock.ExpectRollback()
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"CREDIT"}`)
	resp, err := http.Post(url, "application/json", body)
	transaction := model.Transaction{}
	respdata, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respdata, &transaction)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, 500, transaction.Amount)
	assert.EqualValues(t, 1500, transaction.ClosingBalance)
	assert.NoError(t, err)
}

func TestCreateTransactionSuccessForDebit(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockService.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	columns := []string{"ID", "balance"}
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(123, 1000.0))
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(123, 1000.0))
	mockService.Mock.ExpectBegin()
	mockService.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockService.Mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	mockService.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockService.Mock.ExpectCommit()
}

func TestCreateTransactionFailsForWithQueryFailsToInsertTransaction(t *testing.T) {
	mockDatabase := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockDatabase.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	columns := []string{"ID", "Balance"}
	mockDatabase.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 200.0))
	mockDatabase.Mock.ExpectBegin()
	mockDatabase.Mock.ExpectQuery("UPDATE").WillReturnRows()
	mockDatabase.Mock.ExpectQuery("INSERT").WillReturnError(fmt.Errorf("faild to insert"))
	mockDatabase.Mock.ExpectRollback()
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"CREDIT"}`)
	resp, err := http.Post(url, "application/json", body)
	transaction := model.Transaction{}
	respdata, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respdata, &transaction)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, 500, transaction.Amount)
	assert.EqualValues(t, 500, transaction.ClosingBalance)
	assert.NoError(t, err)
}
func TestCreateTransactionSuccessForCREDIT(t *testing.T) {
	mockDatabase := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockDatabase.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	columns := []string{"ID", "balance"}
	mockDatabase.Mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 200.0))
	mockDatabase.Mock.ExpectBegin()
	mockDatabase.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectCommit()
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"CREDIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, err)
}
func TestCreateTransactionSuccessForDEBIT(t *testing.T) {
	mockDatabase := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mockDatabase.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	columns := []string{"ID", "balance"}
	mockDatabase.Mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 200))
	mockDatabase.Mock.ExpectBegin()
	mockDatabase.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectCommit()
	body := strings.NewReader(`{"wallet_id":123, "amount":120, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, err)
}
func TestRevertTransactionSuccessForDEBIT(t *testing.T) {
	mockDatabase := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction/{tran_id}", mockDatabase.Database, RevertTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction/123"
	transactionColumns := []string{"id", "amount", "type"}
	walletColumns := []string{"id", "balance"}
	mockDatabase.Mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows(transactionColumns).AddRow(1, 200.0, "DEBIT"))

	mockDatabase.Mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows(walletColumns).AddRow(1, 200.0))

	mockDatabase.Mock.ExpectBegin()
	mockDatabase.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mockDatabase.Mock.ExpectCommit()
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, err)
}

func TestCreateDebitTransaction(t *testing.T) {
	mock := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction", mock.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction"
	body := strings.NewReader(`{"wallet_id":123, "amount":500, "type":"DEBIT"}`)
	resp, err := http.Post(url, "application/json", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, err)
}

func TestRevertTransactionSuccess(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/transaction/{tran_id}", mockService.Database, CreateTransaction)
	defer testService.Server.Close()
	url := testService.Server.URL + "/transaction/1"
	columns := []string{"ID", "amount", "wallet_id"}
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 1000.0, 1))
	// mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(123, 1000.0))
	// mockService.Mock.ExpectBegin()
	// mockService.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	// mockService.Mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	// mockService.Mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	// mockService.Mock.ExpectCommit()
	// client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	transaction := model.Transaction{}
	respdata, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respdata, &transaction)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// assert.EqualValues(t, 500, transaction.Amount)
	// assert.EqualValues(t, 1500, transaction.ClosingBalance)
	// assert.NoError(t, err)
}
