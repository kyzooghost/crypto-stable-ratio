package main

import (
	"encoding/csv"
	"fmt"
	"os"
	// "github.com/aws/aws-lambda-go/lambda"
)

// type MyEvent struct {
// 	Name string `json:"name"`
// }

// func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
// 	return fmt.Sprintf("Hello %s!", name.Name), nil
// }

func main() {
	seedValues, err := parseLocalSeedValues()
	for i := range seedValues {
		fmt.Printf("%s\n", seedValues[i])
	}
	// lambda.Start(HandleRequest)
}

func parseLocalSeedValues() ([][]string, error) {
	file, err := os.Open("./seed-values.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %s", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV records: %s", err)
	}

	return records, nil
}

// Parse seed-values.csv into memory
// Query API endpoints
// Merge API endpoint data, with local seed-values.csv
// Save result to S3 bucket
