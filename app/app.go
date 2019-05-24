package app

import (
	"fmt"
	"log"
	"net/http"
	"wallet/app/handler"

	"wallet/app/model"
	"wallet/config"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/urfave/negroni"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

type Route struct {
	route   string
	handler func(w http.ResponseWriter, r *http.Request)
	method  string
}

func (a *App) InitializeAndRun(config *config.Config, port string) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal(fmt.Sprintf("connection failed to dbwith err : %#v  ", err.Error()))
	}
	a.DB = model.DBMigrate(db)
	router := mux.NewRouter()
	routes := getRouter(a)
	for _, route := range routes {
		router.HandleFunc(route.route, route.handler).Methods(route.method)
	}
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(router)
	http.ListenAndServe(port, n)
}

func getRouter(a *App) []Route {
	return []Route{
		{
			route:   "/walletapi/wallet/{wallet_id}",
			handler: a.GetWallet(),
			method:  "GET",
		},
		{
			route:   "/walletapi/wallet",
			handler: a.CreateWallet(),
			method:  "POST",
		},
		{
			route:   "/walletapi/wallet/{wallet_id}/transactions",
			handler: a.GetWalletTransactions(),
			method:  "GET",
		},
		{
			route:   "/walletapi/transaction",
			handler: a.CreateTransaction(),
			method:  "POST",
		},
		{
			route:   "/walletapi/transaction/{tran_id}",
			handler: a.RevertTransaction(),
			method:  "DELETE",
		},
	}
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

//toDo move this all wrapper to Handler itself

func (a *App) GetWallet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.GetWallet(a.DB, w, r)
	}
}
func (a *App) CreateWallet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.CreateWallet(a.DB, w)
	}
}

func (a *App) GetWalletTransactions() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.GetWalletTransactions(a.DB, w, r)
	}
}

func (a *App) CreateTransaction() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.CreateTransaction(a.DB, w, r)
	}
}

func (a *App) RevertTransaction() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.RevertTransaction(a.DB, w, r)
	}
}
