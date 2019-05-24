package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"wallet/app/constant"
	"wallet/app/model"
)

func CreateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&transaction); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	wallet := getWalletFor(db, transaction.WalletId)
	if err := db.Error; err != nil {
		respondError(w, http.StatusInternalServerError, "failed while fetching wallet information")
	}
	if !isValidTransactionType(transaction) {
		respondError(w, http.StatusBadRequest, "invalid transaction type")
		return
	}
	tran, err := processTransaction(wallet, transaction, db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to process transaction, "+err.Error())
		return
	}
	respondSuccess(w, *tran)
}

func RevertTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tranId, err := strconv.ParseInt(vars["tran_id"], 10, 64)
	transaction := model.Transaction{}
	transaction.ID = uint(tranId)
	db.Find(&transaction)

	if err := db.Error; err != nil {
		respondError(w, http.StatusInternalServerError, "failed while transaction information")
	}
	wallet := getWalletFor(db, transaction.WalletId)

	if err := db.Error; err != nil {
		respondError(w, http.StatusInternalServerError, "failed while fetching wallet information")
	}
	transaction = createRevertTransaction(transaction)
	tran, err := processTransaction(wallet, transaction, db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to process transaction, "+err.Error())
		return
	}
	respondSuccess(w, *tran)
}
func rollbackOnError(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	}
}
func processTransaction(wallet model.Wallet, transaction model.Transaction, db *gorm.DB) (*model.Transaction, error) {
	if !canProcessTransaction(transaction, wallet) {
		return nil, fmt.Errorf("cannot process transaction, check your balance")
	}
	tx := db.Begin()
	defer rollbackOnError(tx)
	if err := tx.Error; err != nil {
		return nil, err
	}
	wallet.Balance = getUpdatedWalletBalance(wallet, transaction)
	if err := tx.Save(&wallet).Error; err != nil {
		return nil, err
	}
	transaction.ClosingBalance = wallet.Balance
	if err := tx.Save(&transaction).Error; err != nil {
		return nil, err
	}
	if err := tx.Save(&wallet).Error; err != nil {
		return nil, err
	}
	return &transaction, tx.Commit().Error
}

func getUpdatedWalletBalance(wallet model.Wallet, transaction model.Transaction) float32 {
	if transaction.Type == constant.CREDIT {
		return wallet.Balance + transaction.Amount
	}
	return wallet.Balance - transaction.Amount
}

func canProcessTransaction(transaction model.Transaction, wallet model.Wallet) bool {
	return constant.CREDIT == transaction.Type ||
		wallet.Balance-transaction.Amount >= 0
}

func isValidTransactionType(transaction model.Transaction) bool {
	return constant.CREDIT == transaction.Type || constant.DEBIT == transaction.Type
}

func getWalletFor(db *gorm.DB, walletId uint) model.Wallet {
	wallet := model.Wallet{}
	wallet.ID = walletId
	db.First(&wallet)
	return wallet
}

func createRevertTransaction(transaction model.Transaction) model.Transaction {
	updatedTran := model.Transaction{}
	if transaction.Type == constant.CREDIT {
		updatedTran.Type = constant.DEBIT
	} else {
		updatedTran.Type = constant.CREDIT
	}
	updatedTran.Amount = transaction.Amount
	updatedTran.Description = fmt.Sprint("Revert of :", transaction.ID)
	updatedTran.WalletId = transaction.WalletId
	return updatedTran
}
