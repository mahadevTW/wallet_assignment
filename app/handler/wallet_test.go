package handler

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"wallet/app/model"
	"wallet/testutils"
)

func TestGetWalletFailsForNonNumericWalletId(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/wallet/{wallet_id}", mockService.Database, GetWallet)
	defer testService.Server.Close()
	url := testService.Server.URL + "/wallet/abc"
	resp, err := http.Get(url)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)
}

func TestGetWalletFailsForInvalidWalletId(t *testing.T) {
}

func TestGetWalletSuccess(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/wallet/{wallet_id}", mockService.Database, GetWallet)
	defer testService.Server.Close()
	url := testService.Server.URL + "/wallet/123"
	columns := []string{"ID", "balance"}
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(123, 400))

	resp, err := http.Get(url)
	wallet := model.Wallet{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &wallet)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, 123, wallet.ID)
	assert.EqualValues(t, 400, wallet.Balance)
	assert.NoError(t, err)
}

func TestGetWalletTransactionsFailsForInvalidTransaction(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/wallet/{wallet_id}/transactions", mockService.Database, GetWalletTransactions)
	defer testService.Server.Close()
	url := testService.Server.URL + "/wallet/abc/transactions"
	resp, err := http.Get(url)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)
}

func TestGetWalletTransactionsFailsDBQueryFailsForGetWallet(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/wallet/{wallet_id}", mockService.Database, GetWallet)
	defer testService.Server.Close()
	url := testService.Server.URL + "/wallet/123"
	columns := []string{"ID", "balance"}
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns))
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)
}

func TestGetWalletTransactionsFailsDBQueryFailsForGetTransactions(t *testing.T) {

}

func TestGetWalletTransactionsSuccess(t *testing.T) {

}
