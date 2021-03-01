package auth

import (
	"net/http"

	"github.com/harriklein/pBE/pBEServer/app"
)

// Init initializes the endpoint
func Init() {
	// Initialize handler and register routers
	_authRouter := app.SrvMux.PathPrefix("/api/v1").Subrouter()
	_authRouter.HandleFunc("/signup", handleSignUp).Methods(http.MethodGet)
	_authRouter.HandleFunc("/login", handleLoginBasic).Methods(http.MethodGet)
	_authRouter.HandleFunc("/logout", handleLogout).Methods(http.MethodGet)

	_authTestRouter := app.SrvMux.PathPrefix("/api/v1/login").Subrouter()
	_authTestRouter.HandleFunc("/test", handleLoginTest).Methods(http.MethodGet)
	_authTestRouter.Use(CheckJWT)
}
