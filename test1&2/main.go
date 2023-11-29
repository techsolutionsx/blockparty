package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"database/sql"
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
	// Read CIDs from the CSV file
	cids, err := readCIDsFromFile("ipfs_cids.csv")
	if err != nil {
		fmt.Println("Error reading CIDs:", err)
		return
	}

	// Iterate through the CIDs and fetch metadata
	for _, cid := range cids {
		url := "https://ipfs.io/ipfs/" + cid
		metadata, err := fetchMetadata(url)
		if err != nil {
			fmt.Printf("Error fetching metadata for CID %s: %s\n", cid, err)
			continue
		}

		// Print or process the metadata as needed
		fmt.Printf("Metadata for CID %s:\n%+v\n\n", cid, metadata)

		err = insertMetadataIntoDatabase(metadata)
		if err != nil {
			fmt.Printf("Error inserting metadata into database: %s\n", err)
		}
	}
}

func readCIDsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cids []string
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		cids = append(cids, record[0])
	}

	return cids, nil
}

func fetchMetadata(url string) (*IPFSObject, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch metadata, status code: %d", resp.StatusCode)
	}

	var metadata IPFSObject
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}


func insertMetadataIntoDatabase(metadata *IPFSObject) error {
	// Construct the database connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert metadata into the database
	_, err = db.Exec("INSERT INTO metadata (image, description, name) VALUES ($1, $2, $3)",
		metadata.Image, metadata.Description, metadata.Name)
	if err != nil {
		return err
	}

	return nil
}