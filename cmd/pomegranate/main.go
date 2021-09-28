package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pomegranate/database"
	"pomegranate/newznab"
	"pomegranate/service"
	"pomegranate/themoviedb"
	"syscall"
	"time"
)

const defaultPort = 3000
const themoviedbApiKeyEnvironmentKey = "THEMOVIEDB_API_KEY"
const newznabEnvironmentPrefix = "NEWZNAB"

func loadSettings() (config service.Config, err error) {
	themoviedbApiKey := os.Getenv(themoviedbApiKeyEnvironmentKey)
	if themoviedbApiKey == "" {
		return config, fmt.Errorf("invalid or missing required environment key: %s", themoviedbApiKeyEnvironmentKey)
	}
	config.Tmdb = themoviedb.New(themoviedbApiKey)

	for i := 1; true; i++ {
		newznabHostKey := fmt.Sprintf("%s_HOST_%d", newznabEnvironmentPrefix, i)
		newznabApiKey := fmt.Sprintf("%s_KEY_%d", newznabEnvironmentPrefix, i)

		host := os.Getenv(newznabHostKey)
		apiKey := os.Getenv(newznabApiKey)

		if host == "" {
			break
		}

		config.Newz = append(config.Newz, newznab.Newznab{Host: host, ApiKey: apiKey})
	}
	if len(config.Newz) <= 0 {
		return config, fmt.Errorf("invalid or missing newznab environemnt keys. Use keys %s_HOST_1 and %s_KEY_1 for setting the sources of nzb files. Numbers should be sequential and start at 1. Key is optional if the server does not require one", newznabEnvironmentPrefix, newznabEnvironmentPrefix)
	}

	// TODO: Make the database path a setting
	db, err := database.Open("pomegranate.db")
	if err != nil {
		log.Fatal(fmt.Errorf("database.Open: %w", err))
	}
	config.DB = db

	return
}

func main() {
	fmt.Println("Pomegranate is initializing...")

	config, err := loadSettings()
	if err != nil {
		log.Fatal(fmt.Errorf("loadSettings: %w", err))
	}

	dbKeys, err := config.DB.BucketKeys(database.MovieBucketName)
	if err != nil {
		log.Fatal(fmt.Errorf("db.BucketKeys: %w", err))
	}

	fmt.Printf("Movies in database (%d):", len(dbKeys))
	for _, k := range dbKeys {
		fmt.Printf(" %s", string(k))
	}
	fmt.Printf("\n")

	addr := fmt.Sprintf(":%d", defaultPort)
	server := &http.Server{Addr: addr, Handler: service.Service(config)}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	signalListener(server, serverCtx, serverStopCtx)

	fmt.Printf("Listening on %s\n", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func signalListener(server *http.Server, serverCtx context.Context, serverStopCtx context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
}
