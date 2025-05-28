package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

// Response represents the complete block mode response structure
type Response struct {
	WorkflowRunID string `json:"workflow_run_id"`
	TaskID        string `json:"task_id"`
	Data          struct {
		ID          string                 `json:"id"`
		WorkflowID  string                 `json:"workflow_id"`
		Status      string                 `json:"status"`
		Outputs     map[string]interface{} `json:"outputs"`
		Error       interface{}            `json:"error"`
		ElapsedTime float64                `json:"elapsed_time"`
		TotalTokens int                    `json:"total_tokens"`
		TotalSteps  int                    `json:"total_steps"`
		CreatedAt   int64                  `json:"created_at"`
		FinishedAt  int64                  `json:"finished_at"`
	} `json:"data"`
}

type WorkflowRunRequest struct {
	Inputs       map[string]interface{} `json:"inputs"`
	ResponseMode string                 `json:"response_mode"` // "streaming" or "blocking"
	User         string                 `json:"user"`          // user identifier
}

func RunDifyWorkflow(request WorkflowRunRequest) (Response, error) {
	url := viper.GetString("dify.api_endpoint") + "/workflows/run"

	payload, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		panic(fmt.Errorf("failed to create request: %w", err))
	}

	req.Header.Set("Authorization", "Bearer "+viper.GetString("dify.api_key"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Response{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check for errors
	if response.Data.Error != nil {
		return response, fmt.Errorf("error in response: %v", response.Data.Error)
	}

	// Return the response
	return response, nil

}
