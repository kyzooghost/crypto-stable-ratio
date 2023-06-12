package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	// "github.com/aws/aws-lambda-go/lambda"
)

const DAY = 86400
const WEEK = 7 * DAY

// type MyEvent struct {
// 	Name string `json:"name"`
// }

// func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
// 	return fmt.Sprintf("Hello %s!", name.Name), nil
// }

func main() {
	// seedValues, _ := parseLocalSeedValues()
	// for i := range seedValues {
	// 	fmt.Printf("%s\n", seedValues[i])
	// }
	// lambda.Start(HandleRequest)
	_getSanitizedData(os.Getenv("STABLE_MCAP_ENDPOINT"))
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

type Response struct {
	Success bool        `json:"success"`
	Data    []Datapoint `json:"data"`
}

type Datapoint struct {
	Timestamp int64   `json:"t"`
	Value     float64 `json:"v"`
}

func _getSanitizedData(uri string) ([]Datapoint, error) {
	data := make([]Datapoint, 0)

	// Input validation
	if uri == "" {
		println("Did not provide request URI")
		return data, fmt.Errorf("Did not provide request URI")
	}

	// Make HTTP GET request
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(uri)
	if err != nil {
		println("_getRequest GET request error", err)
		return data, fmt.Errorf("_getRequest GET request error: %s", err)
	}
	defer resp.Body.Close()

	// Parse response
	var parsedResponse Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResponse)
	if err != nil {
		println("_getRequest decode response error", err)
		return data, fmt.Errorf("_getRequest decode response error: %s", err)
	}
	data = parsedResponse.Data

	// Sanitize data
	sanitizedData := make([]Datapoint, 0)
	var tempTimestamp int64
	var tempData Datapoint
	var tempLength int

	for i := range data {
		tempData = data[i]
		// Truncate timestamp from ms to s
		tempTimestamp = tempData.Timestamp / 1000
		// Move up by 4 days (So that Sunday - Sat in the same week, will round down to the same timestamp)
		tempTimestamp = tempTimestamp + 4*DAY
		// Unix timestamp start on Thursday 00:00 UTC, so if we didn't add days then Wednesday and Friday in the same week would round down to different timestamp
		tempTimestamp = tempTimestamp / WEEK * WEEK
		tempTimestamp = tempTimestamp - 4*DAY
		tempData.Timestamp = tempTimestamp
		// Result in timestamp rounded down to most recent Sunday 00:00 UTC
		data[i].Timestamp = tempData.Timestamp

		// We only store most recent result for each week
		// We ASSUME endpoint providing data in earliest-to-latest order
		tempLength = len(sanitizedData)
		// We do not have a datapoint for current Sunday 00:00 timestamp being considered
		if tempLength == 0 || sanitizedData[tempLength-1].Timestamp < tempTimestamp {
			sanitizedData = append(sanitizedData, tempData)
		} else {
			// ASSUME the endpoint provided data in earliest-to-latest order, so we are dealing with progressively later timestamps
			// Overwrite if encounter later timestamp (which has been rounded down to the same Sunday 00:00)
			sanitizedData[tempLength-1].Value = tempData.Value
		}
	}

	fmt.Printf("%v\n", sanitizedData)

	return sanitizedData, nil
}

// Parse seed-values.csv into memory
// Query API endpoints
// Merge API endpoint data, with local seed-values.csv
// Save result to S3 bucket
