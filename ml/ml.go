package ml

import (
	"backend/cache"
	"bufio"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	mlAddress            string
	ErrMladdressNotFound = errors.New("ml address not found")
	retry                = 5
)

type Config struct {
	address string
}

func GetConfig() (*Config, error) {

	cfg := &Config{}

	mladdrEnv := os.Getenv("ML_ADDRESS")
	if mladdrEnv != "" {
		cfg.address = mladdrEnv
	} else {
		return nil, ErrMladdressNotFound
	}

	return cfg, nil
}

func SetUp(cfg *Config) error {
	slog.Debug("Starting ml...")

	mlAddress = "http://" + cfg.address

	_, err := http.Get(mlAddress)
	times := retry
	for err != nil && times > 0 {
		_, err = http.Get(mlAddress)
		times--
		time.After(time.Second)
	}
	if err != nil {
		return err
	}

	slog.Debug("ml started...")

	return nil
}

func GetRecomendation(rawJson []byte) (string, error) {
	slog.Debug("Getting recomendation...")

	recom, err := cache.Client.Get(context.Background(), string(rawJson)).Result()
	if err != nil {
		slog.Error(err.Error())
	}

	if recom == "" {
		slog.Debug("Empty cache")
		resb, err := GetRecomendationFromModel(rawJson)
		if err != nil {
			return "", err
		}
		return resb, nil
	}

	return recom, nil
}

func GetRecomendationFromModel(rawJson []byte) (string, error) {
	slog.Debug("Getting recomendation from model...")

	buf := bytes.NewReader(rawJson)
	res, err := http.Post(mlAddress, "application/json", buf)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	jsonstr := ""
	for scanner.Scan() {
		jsonstr += scanner.Text()
	}

	//caching
	slog.Debug("Caching...")
	err = cache.Client.SetEx(context.Background(), string(rawJson), jsonstr, time.Minute*10).Err()
	if err != nil {
		slog.Error(err.Error())
	}

	return jsonstr, nil
}
