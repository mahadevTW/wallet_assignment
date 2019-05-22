package testutils

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Mock struct {
	DB       *sql.DB
	Mock     sqlmock.Sqlmock
	Database *gorm.DB
}
type TestServer struct {
	request *mux.Router
	Server  *httptest.Server
}

func NewMockDb(t *testing.T) *Mock {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to mock databnase")
	database, _ := gorm.Open("mysql", db)
	database.LogMode(true)
	return &Mock{db, mock, database}
}
func (m *Mock) ExpectQuery(query string) *Mock {
	m.Mock.ExpectQuery(query)
	return m
}

func (m *Mock) ExpectQueryWithResult(query string, rows *sqlmock.Rows) *Mock {
	m.Mock.ExpectQuery(query).WillReturnRows(rows)
	return m
}

func NewTestServer() *TestServer {
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	return &TestServer{r, ts}
}

func (t *TestServer) RegisterHandler(path string, db *gorm.DB, handlerFunc func(db *gorm.DB, w http.ResponseWriter, r *http.Request)) *TestServer {
	t.request.HandleFunc(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(db, w, r)
	}))
	return t
}
