package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	newRelicAPIURL = "https://insights-api.newrelic.com/v1/accounts/6225327/query"
	query          = "SELECT rate(count(*), 1 minute) FROM Transaction WHERE appName = 'devfest-app-2024' SINCE 5 minutes ago"
	newRelicAPIKey = "0bc1f33be4ec7b4cafa463143a48bb09FFFFNRAL" // Replace with your API key
)

// Response structure to parse New Relic API response
type NewRelicResponse struct {
	Results []map[string]interface{} `json:"results"`
}

func getNewRelicMetrics() (float64, error) {
	req, err := http.NewRequest("POST", newRelicAPIURL, bytes.NewBuffer([]byte(fmt.Sprintf("query=%s", query))))
	if err != nil {
		return 0, err
	}

	// Set API Key for Authorization
	req.Header.Set("X-Query-Key", newRelicAPIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result NewRelicResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	// Assuming we have one result with the rate in it
	if len(result.Results) > 0 {
		if val, ok := result.Results[0]["rate(count(*), 1 minute)"].(float64); ok {
			return val, nil
		}
	}
	return 0, fmt.Errorf("unexpected response format")
}

func main() {
	metricValue, err := getNewRelicMetrics()
	if err != nil {
		log.Fatalf("Error fetching metrics from New Relic: %v", err)
	}

	// Now, use the fetched metric to make a prediction using Vertex AI
	fmt.Printf("Fetched New Relic metric: %f\n", metricValue)

	// Call your Vertex AI prediction function here
	// For example: MakePrediction(metricValue)
}
