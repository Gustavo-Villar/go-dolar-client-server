package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type USD2BRL struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

var db *sql.DB

func init() {
	var err error
	// Initialize the SQLite database
	db, err = sql.Open("sqlite", "./quotations.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS quotations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		createDate TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func getQuotationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Route not found"})
		return
	}

	// Fetch the quotation from the API
	data, err := getQuotation()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error fetching quotation"})
		return
	}

	// Save the quotation to the database
	err = saveQuotation(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error saving quotation"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func getQuotation() (*USD2BRL, error) {
	// Create a context with a timeout of 200ms
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Request URL
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	// Create an HTTP request with the context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Make the HTTP GET request to the AwesomeAPI endpoint
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error is due to the context timeout
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request to AwesomeAPI timed out")
		}
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal the JSON response into the USD2BRL struct
	var quotation USD2BRL
	err = json.NewDecoder(resp.Body).Decode(&quotation)
	if err != nil {
		return nil, err
	}

	return &quotation, nil
}

func saveQuotation(data *USD2BRL) error {
	// Create a context with a timeout of 10ms
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Execute the insert statement within the context
	_, err := db.ExecContext(ctx, `INSERT INTO quotations (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		data.USDBRL.Code, data.USDBRL.Codein, data.USDBRL.Name, data.USDBRL.High, data.USDBRL.Low, data.USDBRL.VarBid, data.USDBRL.PctChange, data.USDBRL.Bid, data.USDBRL.Ask, data.USDBRL.Timestamp, data.USDBRL.CreateDate)
	if err != nil {
		// Check if the error is due to the context timeout
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Database operation timed out")
		}
		return err
	}

	return nil
}

func main() {
	http.HandleFunc("/", getQuotationHandler)

	fmt.Println("Server Running on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
