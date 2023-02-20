package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"

	"github.com/PSNAppz/Fold-ELK/db"
	"github.com/PSNAppz/Fold-ELK/handler"
	"github.com/rs/zerolog"
)

func SetUpRouter() *gin.Engine {
	var err error
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	// Set default config for testing
	dbConfig := db.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "fold-elk",
		Password: "password",
		DbName:   "fold_elk",
		Logger:   logger,
	}
	logger.Info().Interface("config", &dbConfig).Msg("config:")
	dbInstance, err := db.Init(dbConfig)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	logger.Info().Msg("Database connection established")

	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}

	h := handler.New(dbInstance, esClient, logger)
	router := gin.Default()

	rg := router.Group("/v1")
	h.Register(rg)

	return router
}

func TestUserHandler(t *testing.T) {
	router := SetUpRouter()
	router.Group("/v1")

	// Test CREATE operation
	user := map[string]string{
		"name": "PSN",
	}
	payload, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	//get user id from resp
	var result map[string]interface{}
	json.Unmarshal([]byte(resp.Body.Bytes()), &result)
	userID := int(result["user"].(map[string]interface{})["ID"].(float64))

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.Code)
	}

	// Test READ operation
	req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/users/%d", userID), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Test UPDATE operation
	update := map[string]string{
		"name": "NEW",
	}
	payload, _ = json.Marshal(update)
	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/v1/users/%d", userID), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Test DELETE operation
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/users/%d", userID), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 204, got %d", resp.Code)
	}
}
