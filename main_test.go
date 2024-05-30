package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	API_URL = os.Getenv("API_URL")
	if API_URL == "" {
		log.Fatalf("API_URL is not set in the .env file")
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Pair",
			query:          "pair=XXBTZUSD",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ltp":[{"pair":"XXBTZUSD"`,
		},
		{
			name:           "Multiple Valid Pairs",
			query:          "pair=XXBTZUSD&pair=XXBTZEUR",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ltp":[{"pair":`,
		},
		{
			name:           "Invalid Pair",
			query:          "pair=INVALIDPAIR",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `pair not found`,
		},
		{
			name:           "No Pair Specified",
			query:          "",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ltp":[{"pair":`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/ltp?"+tt.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Handler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
