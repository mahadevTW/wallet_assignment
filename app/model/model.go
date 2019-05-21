package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Wallet struct {
	gorm.Model
	Balance float32
}
type Transaction struct {
	gorm.Model
	Amount         float32 `json:"amount"`
	Type           string  `gorm:"type:ENUM('CREDIT','DEBIT');" json:"type"`
	ClosingBalance float32
	Description    string
	WalletId       uint `json:"wallet_id"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.LogMode(true)
	db.AutoMigrate(&Wallet{}, &Transaction{})
	db.Model(&Transaction{}).AddForeignKey("wallet_id", "wallets(id)", "CASCADE", "CASCADE")
	return db
}
