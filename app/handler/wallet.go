package handler

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"wallet/app/model"
)

func CreateWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	wallet := model.Wallet{}
	db.Create(&wallet)
	respondSuccess(w, wallet)
}

func GetWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletId, err := strconv.ParseInt(vars["wallet_id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid wallet id")
		return
	}
	wallet := model.Wallet{}
	wallet.ID = uint(walletId)
	db.First(&wallet)
	respondSuccess(w, wallet)
}

func GetWalletTransactions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletId, err := strconv.ParseInt(vars["wallet_id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid wallet id")
		return
	}
	var transactions []model.Transaction
	wallet := model.Wallet{}
	wallet.ID = uint(walletId)
	db.First(&wallet)
	db.Where("wallet_id=?", walletId).Find(&transactions)
	respondSuccess(w, transactions)
}
