package server

import (
	"backend/ml"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

/*
	{
	    "clientId":"client", // ИД пользователя
	    "organizationId":"organization", // ИД организациии
	    "segment":"value", // Сегмент организации: "Малый бизнес", "Средний бизнес", "Крупный бизнес"
	    "role":"value", // Роль уполномоченного лица: "ЕИО", "Сотрудник"
	    "organizations": 3, // Общее количество организаций у уполномоченного лица: 1..300
	    "currentMethod": "method", // Действующий способ подписания."SMS", "PayControl", "КЭП на токене", "КЭП в приложении"
	    "mobileApp": true, // Наличие мобильного приложения
	    "signatures": { // Подписанные ранее типы документов
	        "common": {
	            "mobile":3, // Количество подписанных базовых документов в мобайле
	            "web":10, // Количество подписанных базовых документов в вебе
	        },
	        "special": {
	            "mobile":5, // Количество подписанных документов особой важности в мобайле
	            "web":6, // Количество подписанных документов особой важности в вебе
	        }
	    },
	    "availableMethods":["method1", "method2"], // Уже подключенные способы подписания."SMS", "PayControl", "КЭП на токене", "КЭП в приложении"
	    "claims": 0 // Наличие обращений в банк по причине проблем с использованием СМС
	}
*/
type ClientInfo struct {
	ClientId         string              `json:"clientId"`
	OrganizationId   string              `json:"organizationId"`
	Segment          string              `json:"segment"`
	Role             string              `json:"role"`
	Organizations    int                 `json:"organizations"`
	CurrentMethod    string              `json:"currentMethod"`
	MobileApp        bool                `json:"mobileApp"`
	Signatures       map[string]SignInfo `json:"signatures"`
	AvailableMethods []string            `json:"availableMethods"`
	Claims           int                 `json:"claims"`
	Context          string              `json:"context"`
}

type SignInfo struct {
	Mobile int `json:"mobile"`
	Web    int `json:"web"`
}

func setupHandlers() {
	http.HandleFunc("/recomendation/{$}", recomendation)
}

func recomendation(w http.ResponseWriter, r *http.Request) {

	ci := &ClientInfo{}

	rawBody := make([]byte, r.ContentLength)
	n, err := io.ReadFull(r.Body, rawBody)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if n == 0 {
		slog.Error("Empty body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(rawBody, ci)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//data validation
	if ci.CurrentMethod != "SMS" {
		slog.Debug("Method is not SMS")
		w.WriteHeader(http.StatusOK)
		return
	}
	slog.Debug("Method is SMS")

	recomendation, err := ml.GetRecomendation(rawBody)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(recomendation))
}
