package main

import (
	"backend/cache"
	"backend/ml"
	"backend/server"
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
}

func init() {

	devEnv := flag.Bool("dev", false, "Load development enviroment variables (for testing purposes)")
	flag.Parse()

	if *devEnv {
		err := godotenv.Load()

		lvl := new(slog.LevelVar)
		lvl.Set(slog.LevelDebug)
		logger := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
		slog.SetDefault(slog.New(logger))

		if err != nil {
			slog.Error("Error loading .env file")
			os.Exit(0)
		}
	}

	slog.Info("flag ckeck complited")
}

func main() {
	slog.Info("Starting server...")

	slog.Info("Getting configs...")

	serverCfg := server.GetConfig()
	cacheCfg, err := cache.GetConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}

	mlCfg, err := ml.GetConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}

	slog.Info("Starting modules...")

	err = ml.SetUp(mlCfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}

	err = cache.Start(cacheCfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}

	err = server.Start(serverCfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}

}
