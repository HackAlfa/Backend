package server

import (
	"log/slog"
	"net/http"
	"os"
)

var server *http.Server

type Config struct {
	port string
}

func GetConfig() *Config {

	port := "8080"
	cfg := &Config{}

	portEnv := os.Getenv("SERVER_PORT")
	if portEnv != "" {
		port = portEnv
	}

	cfg.port = port

	return cfg
}

func Start(cfg *Config) error {
	if server != nil {
		slog.Info("Web server already started...")
		return nil
	}

	setupHandlers()

	server = &http.Server{
		Addr: ":" + cfg.port,
	}

	slog.Info("Web server started...")
	return server.ListenAndServe()
}

func Close() error {
	slog.Info("Web server stopped...")
	return server.Close()
}
