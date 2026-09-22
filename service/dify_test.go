package service_test

import (
	"testing"

	"github.com/nbtca/saturday/service"
	"github.com/spf13/viper"
)

func TestRunDifyWorkflow(t *testing.T) {
	if viper.GetString("dify.api_endpoint") == "" {
		t.Skip("dify.api_endpoint is not set")
	}
	// Example request
	request := service.WorkflowRunRequest{
		Inputs: map[string]interface{}{
			"EventId": 49,
		},
		ResponseMode: "blocking", // Change to "blocking" for a single response
		User:         "abc-123",  // User identifier
	}

	response, err := service.RunDifyWorkflow(request)
	if err != nil {
		t.Fatalf("Error running workflow: %v", err)
	}

	t.Logf("Workflow Run ID: %s", response.WorkflowRunID)
	t.Logf("Task ID: %s", response.TaskID)
	t.Logf("Response Data: %+v", response.Data)
}
