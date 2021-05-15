package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/KaiserWerk/SimpleCA/internal/certmaker"
	"github.com/KaiserWerk/SimpleCA/internal/configuration"
	"github.com/KaiserWerk/SimpleCA/internal/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	port = "8880"
)

func main() {
	var err error
	configFilePtr := flag.String("config", "", "The configuration file to use")
	portPtr := flag.String("port", "", "The port to run at")
	useUiPtr := flag.Bool("ui", true, "Adds a simple UI for certificate management")
	flag.Parse()

	if *portPtr != "" {
		port = *portPtr
	}

	if *configFilePtr != "" {
		configuration.SetFileSource(*configFilePtr)
	}
	createdConfig, createdSn, err := configuration.Setup()
	if err != nil {
		log.Fatalf("could not set up configuration: %s", err.Error())
	}

	if createdConfig {
		log.Printf("The configuration file was not found; created\n\tStop execution? (y,n)")
		var answer string
		_, _ = fmt.Scanln(&answer)
		if answer == "y" {
			log.Fatalf("Okay, stopped.")
		}
	}

	if createdSn {
		log.Printf("The serial number file not found; created")
	}

	// create root cert and key, if non-existent
	err = certmaker.SetupCA()
	if err != nil {
		log.Fatalf("could not set up CA: %s", err.Error())
	}

	host := fmt.Sprintf(":%s", port)
	router := mux.NewRouter()
	setupRoutes(router, *useUiPtr)

	notify := make(chan os.Signal)
	signal.Notify(notify, os.Interrupt)

	srv := &http.Server{
		Addr:              host,
		Handler:           router,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       3 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		<-notify
		log.Println("Initiating graceful shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
		defer cancel()
		// do necessary stuff here before we exit

		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Fatal("Could not gracefully shut down server: " + err.Error())
		}
	}()

	log.Printf("Server listening on %s...\n", host)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v\n", err.Error())
	}
	log.Println("Server shutdown complete. Have a nice day!")
}

func setupRoutes(router *mux.Router, ui bool) {
	if ui {
		router.HandleFunc("/", handler.IndexHandler).Methods(http.MethodGet)
		router.HandleFunc("/login", handler.LoginHandler).Methods(http.MethodGet, http.MethodPost)
		router.HandleFunc("/logout", handler.LogoutHandler).Methods(http.MethodGet)
		router.HandleFunc("/add", handler.AddCertificateHandler).Methods(http.MethodGet, http.MethodPost)
		router.HandleFunc("/revoke", handler.RevokeCertificateHandler).Methods(http.MethodGet, http.MethodPost) // TODO implement
	}
	router.HandleFunc("/api/certificate/request", handler.ApiRequestCertificateHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/certificate/{id}/obtain", handler.ApiObtainCertificateHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/privatekey/{id}/obtain", handler.ApiObtainPrivateKeyHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/ocsp/", handler.ApiOcspRequestHandler).Methods(http.MethodPost) // TODO only post?
}

