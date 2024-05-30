package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var API_URL string

type RawResponse struct {
	Result map[string]Pair `json:"result"`
}

type Pair struct {
	Price []string `json:"a"`
}

type LTPResponse struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type Response struct {
	LTP []LTPResponse `json:"ltp"`
}

/*
Handler handle the user request

if we receibe a query param named pair, we only return that pair
otherwise returns all the set of pairs
*/
func Handler(w http.ResponseWriter, r *http.Request) {
	pairs := r.URL.Query()["pair"]
	if len(pairs) == 0 {
		pairs = []string{"XXBTZUSD", "XBTCHF", "XXBTZEUR"}
	}

	var wg sync.WaitGroup
	results := make(chan LTPResponse, len(pairs))
	errors := make(chan error, len(pairs))

	for _, pair := range pairs {
		wg.Add(1)
		go getPairPrice(pair, &wg, results, errors)
	}

	wg.Wait()
	close(results)
	close(errors)

	var ltps []LTPResponse
	for res := range results {
		ltps = append(ltps, res)
	}

	for err := range errors {
		if err != nil {
			log.Println("Error fetching LTP:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := Response{LTP: ltps}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*
getPairPrice go to get the pair actual price to the external API

this function use goroutine to optimize the request

for this example, we only get 3 pais, but if we want to add more pairs,
the time request will be the same, cause the go routines permit to make
multiple pair calls at the same time
*/
func getPairPrice(pair string, wg *sync.WaitGroup, results chan<- LTPResponse, errors chan<- error) {
	defer wg.Done()
	resp, err := http.Get(API_URL + pair)
	fmt.Printf("Sending request to %s%s...", API_URL, pair)
	if err != nil {
		log.Println("Error sending request to external API:", err)
		errors <- err
		return
	}
	defer resp.Body.Close()

	var rr RawResponse
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		log.Println("Error decoding the response:", err)
		errors <- err
		return
	}

	rawPair, ok := rr.Result[pair]
	if !ok {
		log.Println("Error getting the pair:", err)
		errors <- fmt.Errorf("pair not found")
		return
	}

	amount, err := strconv.ParseFloat(rawPair.Price[0], 64)
	if err != nil {
		log.Println("Error converting the pair price:", err)
		errors <- err
		return
	}

	results <- LTPResponse{Pair: pair, Amount: amount}
}

func main() {
	//.env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//env variables
	API_URL = os.Getenv("API_URL")
	if API_URL == "" {
		log.Fatalf("API_URL is not set in the .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not set in the .env file")
	}

	//run server
	http.HandleFunc("/api/v1/ltp", Handler)
	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not listen on %s : %v\n", port, err)
	}
}
