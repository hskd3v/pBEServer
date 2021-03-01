package docs

import (
	"net/http"

	"github.com/harriklein/pBE/pBEServer/app"
)

// Init initializes the endpoint
func Init() {
	// region DOCUMENTATION HANDLER ----------------------------
	_docHandler := NewDocHandler()
	app.SrvMux.HandleFunc("/docs/{file}", _docHandler.Get).Methods(http.MethodGet)
	app.SrvMux.HandleFunc("/docs/", _docHandler.Get).Methods(http.MethodGet)
	app.SrvMux.HandleFunc("/docs", _docHandler.Get).Methods(http.MethodGet)
	// endregion -----------------------------------------------
}
