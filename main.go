package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getPrediction(inputData map[string]interface{}) (map[string]interface{}, error) {
	vertexEndpoint := "https://us-central1-aiplatform.googleapis.com/v1/projects/PROJECT_ID/locations/us-central1/endpoints/ENDPOINT_ID:predict"
	token := os.Getenv("GOOGLE_ACCESS_TOKEN")

	jsonData, err := json.Marshal(inputData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", vertexEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func main() {
	// Replace this with real input data
	inputData := map[string]interface{}{
		"instances": []map[string]interface{}{
			{
				"input_data_field": "example",
			},
		},
	}

	prediction, err := getPrediction(inputData)
	if err != nil {
		fmt.Println("Error fetching prediction:", err)
		return
	}
	fmt.Println("Prediction result:", prediction)
}
