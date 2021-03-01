package dummy

import (
	"net/http"

	"github.com/harriklein/pBE/pBEServer/app"
)

var (
	dbRead         = mysqlRead
	dbCreate       = mysqlCreate
	dbUpdate       = mysqlUpate
	dbDelete       = mysqlDelete
	dbApplyUpdates = mysqlApplyUpdates
)

// Init initializes the endpoint
func Init() {
	// Initialize handler and register routers
	_dummyRouter := app.SrvMux.PathPrefix("/api/v1/dummies").Subrouter()
	_dummyRouter.HandleFunc("", handleRead).Methods(http.MethodGet)
	_dummyRouter.HandleFunc("", handleCreateOrApplyUpdates).Methods(http.MethodPost)
	_dummyRouter.HandleFunc("", handleUpdate).Methods(http.MethodPut)
	_dummyRouter.HandleFunc("", handleDelete).Methods(http.MethodDelete)
	_dummyRouter.HandleFunc("/{id}", handleRead).Methods(http.MethodGet)
	_dummyRouter.HandleFunc("/{id}", handleCreateOrApplyUpdates).Methods(http.MethodPost)
	_dummyRouter.HandleFunc("/{id}", handleUpdate).Methods(http.MethodPut)
	_dummyRouter.HandleFunc("/{id}", handleDelete).Methods(http.MethodDelete)
}
