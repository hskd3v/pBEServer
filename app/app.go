package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/harriklein/pBE/pBEServer/log"
)

var (
	//SrvMux variable
	SrvMux  *mux.Router
	srvHTTP *http.Server
)

// Init initializes the main variables of the application
func Init() {

	// Initialize a new server mux -----------------------------
	SrvMux = mux.NewRouter()

}

// RunServer runs the HTTP server
func RunServer() {

	// create a new server
	srvHTTP = &http.Server{
		Addr: cfgAppBindAddress,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      SrvMux, // Pass our instance of gorilla/mux in.
	}

	// start the server
	go func() {
		log.Log.Printf("Starting server on %s", srvHTTP.Addr)

		_error := srvHTTP.ListenAndServe()
		if _error != nil {
			log.Log.Errorf("Error starting server: %s\n", _error)
			os.Exit(1)
		}
	}()

}

// Finish finalizes in the graceful way the application
func Finish() {

	_c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(_c, os.Interrupt)
	signal.Notify(_c, os.Kill)

	// Block until we receive our signal.
	_signal := <-_c
	log.Log.Println("Got signal: ", _signal)

	// Create a deadline to wait for.
	_ctx, _cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer _cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srvHTTP.Shutdown(_ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Log.Println("Shutting down")
	os.Exit(0)

}
