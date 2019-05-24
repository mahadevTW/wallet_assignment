package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet/app/model"
	"wallet/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetWalletFailsForNonNumericWalletId(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/wallet/{wallet_id}", mockService.Database, GetWallet)
	defer testService.Server.Close()
	url := testService.Server.URL + "/wallet/abc"
	resp, err := 	http.Get(url)
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

func TestGetWalletTransactionsTransaction(t *testing.T) {
	mockService := testutils.NewMockDb(t)
	testService := testutils.NewTestServer().RegisterHandler("/wallet/{wallet_id}/transactions", mockService.Database, GetWalletTransactions)
	defer testService.Server.Close()
	url := testService.Server.URL + "/wallet/1/transactions"
	columns := []string{"ID", "amount", "type", "closingbalance", "description"}
	mockService.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 200.0, "CREDIT", 200.0, "credit test"))
	resp, _ := http.Get(url)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	transactions := []model.Transaction{}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &transactions)
	assert.EqualValues(t, "CREDIT", transactions[0].Type)
	assert.EqualValues(t, 200, transactions[0].Amount)
	assert.EqualValues(t, "credit test", transactions[0].Description)
}

func TestGetWalletTransactionsFailsDBQueryFailsForGetTransactions(t *testing.T) {

}

func TestGetWalletTransactionsSuccess(t *testing.T) {

}

func TestCreateWallet(t *testing.T) {
	mockparams := testutils.NewMockDb(t)
	writer := httptest.NewRecorder()
	type args struct {
		db *gorm.DB
		w  httptest.ResponseRecorder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"create_wallet_success",
			args{
				mockparams.Database,
				*writer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateWallet(tt.args.db, &tt.args.w)
			if !isEqual(tt.args.w) {
				t.Errorf("Error while creating wallet")
			}
		})
	}
}

func isEqual(res httptest.ResponseRecorder) bool {
	wallet := model.Wallet{}
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &wallet)
	if err != nil {
		return false
	}
	if res.Result().StatusCode != 200 {
		return false
	}
	if wallet.Balance != 0 {
		return false
	}
	return true
}
