package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	dbHost     = "candidate-testing.co2sjmg0hdpm.us-east-2.rds.amazonaws.com"
	dbPort     = 5432
	dbUser     = "jlee"
	dbPassword = "xj98Zs0f7sl2idk3ls"
	dbName     = "jlee"
)

// IPFSObject represents the structure of the metadata to be retrieved
type IPFSObject struct {
	Image       string `json:"image"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

func main() {
	r := mux.NewRouter()

	// Define the API endpoints
	r.HandleFunc("/tokens", getAllTokens).Methods("GET")
	r.HandleFunc("/tokens/{cid}", getTokenByCID).Methods("GET")

	// Start the server
	http.Handle("/", r)
	port := 8080 // Choose the desired port for your API
	fmt.Printf("Server is listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func getAllTokens(w http.ResponseWriter, r *http.Request) {
	// Fetch all data from the database
	data, err := fetchAllDataFromDatabase()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %s", err), http.StatusInternalServerError)
		return
	}

	// Convert data to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getTokenByCID(w http.ResponseWriter, r *http.Request) {
	// Get the CID parameter from the request
	params := mux.Vars(r)
	cid := params["cid"]

	// Fetch data from the database for the specified CID
	data, err := fetchTokenByCIDFromDatabase(cid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data for CID %s: %s", cid, err), http.StatusNotFound)
		return
	}

	// Convert data to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func fetchAllDataFromDatabase() ([]IPFSObject, error) {
	// Construct the database connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query all data from the database
	rows, err := db.Query("SELECT * FROM metadata")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the query result
	var data []IPFSObject
	for rows.Next() {
		var metadata IPFSObject
		err := rows.Scan(&metadata.Image, &metadata.Description, &metadata.Name)
		if err != nil {
			return nil, err
		}
		data = append(data, metadata)
	}

	return data, nil
}

func fetchTokenByCIDFromDatabase(cid string) (*IPFSObject, error) {
	// Construct the database connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query data from the database for the specified CID
	row := db.QueryRow("SELECT * FROM metadata WHERE cid = $1", cid)

	// Process the query result
	var metadata IPFSObject
	err = row.Scan(&metadata.Image, &metadata.Description, &metadata.Name)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}
