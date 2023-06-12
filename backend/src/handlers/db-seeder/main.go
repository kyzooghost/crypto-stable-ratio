// TODO - refactor to separate business logic from handler

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
	// "github.com/aws/aws-lambda-go/lambda"
)

const DAY = 86400
const WEEK = 7 * DAY

type Response struct {
	Success bool        `json:"success"`
	Data    []Datapoint `json:"data"`
}

type Datapoint struct {
	Timestamp int64   `json:"t"`
	Value     float64 `json:"v"`
}

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
	_mergeLocalAndRemoteData()
}

func _mergeLocalAndRemoteData() ([]Datapoint, error) {
	var localData, remoteData, mergedData []Datapoint
	var err error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() error {
		defer wg.Done()
		localData, err = _getLocalSeedValues()
		if err != nil {
			return fmt.Errorf("_getLocalSeedValues error: %s", err)
		}
		return nil
	}()
	go func() error {
		defer wg.Done()
		remoteData, err = _getRemoteCryptoStableData()
		if err != nil {
			return fmt.Errorf("_getRemoteCryptoStableData error: %s", err)
		}
		return nil
	}()
	wg.Wait()

	tempMap := make(map[int64]float64)

	// Prioritize remoteData over localData
	for i := range remoteData {
		tempMap[remoteData[i].Timestamp] = remoteData[i].Value
	}

	for i := range localData {
		_, ok := tempMap[localData[i].Timestamp]
		if !ok {
			tempMap[localData[i].Timestamp] = localData[i].Value
		}
	}

	for key := range tempMap {
		mergedData = append(mergedData, Datapoint{
			Timestamp: key,
			Value:     tempMap[key],
		})
	}

	sort.Slice(mergedData, func(i, j int) bool {
		return mergedData[i].Timestamp < mergedData[j].Timestamp
	})

	fmt.Printf("%v\n", mergedData)

	return mergedData, nil
}

func _getLocalSeedValues() ([]Datapoint, error) {
	file, err := os.Open("./seed-values.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %s", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read past first (header) line of CSV file
	_, err = reader.Read()
	if err != nil {
		println("_getLocalSeedValues read CSV header error", err)
		return nil, fmt.Errorf("_getLocalSeedValues read CSV header error: %s", err)
	}

	// Read all records from the CSV file
	// Basically a while true loop

	var seedValues []Datapoint
	var tempTimestamp int
	var tempValue float64

	for {
		record, err := reader.Read()
		if err != nil {
			// Check for end-of-file error
			if err.Error() == "EOF" {
				break
			}
			println("_getLocalSeedValues read CSV file error", err)
			return nil, fmt.Errorf("_getLocalSeedValues read CSV file error: %s", err)
		}

		tempTimestamp, err = strconv.Atoi(record[0])
		if err != nil {
			println("_getLocalSeedValues read CSV file while parse to int error", err)
			return nil, fmt.Errorf("_getLocalSeedValues read CSV file while parse to int error error: %s", err)
		}

		if record[1] != "" {
			tempValue, err = strconv.ParseFloat(record[1], 64)
			if err != nil {
				println("_getLocalSeedValues read CSV file while parse to float error", err)
				return nil, fmt.Errorf("_getLocalSeedValues read CSV file while parse to float error error: %s", err)
			}
			seedValues = append(seedValues, Datapoint{
				Timestamp: int64(tempTimestamp),
				Value:     tempValue,
			})
		}

	}

	// fmt.Printf("%v\n", seedValues)

	return seedValues, nil
}

func _getRemoteCryptoStableData() ([]Datapoint, error) {
	var stableMcapData, totalMcapData, cryptoStableData []Datapoint
	var err error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() error {
		defer wg.Done()
		stableMcapData, err = _getSanitizedData(os.Getenv("STABLE_MCAP_ENDPOINT"))
		if err != nil {
			return fmt.Errorf("_getSanitizedData for STABLE_MCAP_ENDPOINT error: %s", err)
		}
		return nil
	}()
	go func() error {
		defer wg.Done()
		totalMcapData, err = _getSanitizedData(os.Getenv("TOTAL_MCAP_ENDPOINT"))
		if err != nil {
			return fmt.Errorf("_getSanitizedData for TOTAL_MCAP_ENDPOINT error: %s", err)
		}
		return nil
	}()
	wg.Wait()

	var tempValue float64
	for i := range stableMcapData {
		if stableMcapData[i].Timestamp == totalMcapData[i].Timestamp {
			tempValue = totalMcapData[i].Value / stableMcapData[i].Value

			cryptoStableData = append(cryptoStableData, Datapoint{
				Timestamp: stableMcapData[i].Timestamp,
				Value:     math.Trunc(tempValue*100) / 100,
			})
		}
	}

	// fmt.Printf("%v\n", cryptoStableData)

	return cryptoStableData, nil
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

	return sanitizedData, nil
}

// Parse seed-values.csv into memory
// Query API endpoints
// Merge API endpoint data, with local seed-values.csv
// Save result to S3 bucket
