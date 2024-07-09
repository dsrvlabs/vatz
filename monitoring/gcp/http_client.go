package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPLoggingClient struct {
	client    *http.Client
	apiKey    string
	projectID string
	url       string
}

func NewHTTPLoggingClient(apiKey, projectID string) *HTTPLoggingClient {
	return &HTTPLoggingClient{
		client:    &http.Client{},
		apiKey:    apiKey,
		projectID: projectID,
		url:       fmt.Sprintf("https://logging.googleapis.com/v2/entries:write?key=%s", apiKey),
	}
}

func (c *HTTPLoggingClient) Log(entry LogEntry) error {
	payload := map[string]interface{}{
		"logName": fmt.Sprintf("projects/%s/logs/example-log", c.projectID),
		"resource": map[string]string{
			"type": "global",
		},
		"entries": []map[string]interface{}{
			{
				"jsonPayload": entry.Payload,
				"severity":    entry.Severity.String(),
				"labels":      entry.Labels,
				"timestamp":   entry.Timestamp.Format(time.RFC3339),
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to write log entry: %s", resp.Status)
	}

	return nil
}

func (c *HTTPLoggingClient) Close() error {
	return nil // No resources to close for HTTP client
}
