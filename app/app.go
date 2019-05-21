package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/urfave/negroni"
	"wallet/app/model"
	"wallet/config"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

type Route struct {
	route   string
	handler http.HandlerFunc
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
		log.Fatal("Could not connect database")
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
			route:   "/projects",
			handler: nil,
			method:  "GET",
		},
	}
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
