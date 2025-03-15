// api/matomo.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	matomoClient *MatomoClient
)

// MatomoClient handles communication with the Matomo Tracking API.
type MatomoClient struct {
	baseURL   string
	siteID    int
	authToken string
	client    *http.Client
}

// Event represents a single tracking event for Matomo.
type Event struct {
	ActionName string    `json:"action_name"`
	URL        string    `json:"url"`
	UserAgent  string    `json:"ua"`
	ClientIP   string    `json:"cip"`
	ClientTime time.Time `json:"cdt"` // Use time.Time, format later
}

// InitMatomoClient creates a new Matomo client.
func InitMatomoClient(baseURL string, siteID int, authToken string) {
	matomoClient = &MatomoClient{
		baseURL:   baseURL,
		siteID:    siteID,
		authToken: authToken,
		client:    &http.Client{Timeout: 10 * time.Second}, // Good practice: set a timeout
	}
}

// TrackEvent sends a single CDN access event to Matomo.
func TrackEvent(event Event) error {
	return TrackEvents([]Event{event}) // Reuse TrackEvents for consistency
}

// TrackEvents sends multiple CDN access events to Matomo in bulk.
func TrackEvents(events []Event) error {
	if len(events) == 0 {
		return nil // Nothing to do
	}

	// Build the request parameters.  We use a map for easier construction.
	params := map[string]interface{}{
		"idsite":     matomoClient.siteID,
		"rec":        1, // Required by Matomo
		"token_auth": matomoClient.authToken,
	}

	// Prepare the 'requests' array for bulk tracking.
	requests := make([]string, 0, len(events))
	for _, event := range events {
		eventParams := url.Values{} //Use url.Values for proper encoding
		eventParams.Set("action_name", event.ActionName)
		eventParams.Set("url", event.URL)
		if event.UserAgent != "" {
			eventParams.Set("ua", event.UserAgent)
		}
		if event.ClientIP != "" {
			eventParams.Set("cip", event.ClientIP)
		}
		if !event.ClientTime.IsZero() {
			eventParams.Set("cdt", event.ClientTime.Format("2006-01-02 15:04:05")) // Matomo's expected format
		}
		requests = append(requests, "?"+eventParams.Encode())
	}
	params["requests"] = requests

	// Encode the parameters as JSON.
	requestBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal Matomo request: %w", err)
	}

	// Construct the full URL.
	fullURL := fmt.Sprintf("%s/matomo.php", matomoClient.baseURL) // Standard Matomo endpoint

	// Create the HTTP request.
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create Matomo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json") // Important for bulk tracking

	// Send the request.
	resp, err := matomoClient.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Matomo request: %w", err)
	}
	defer resp.Body.Close() // Crucial: always close the response body

	// Check the response status code.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Matomo API returned non-OK status: %d", resp.StatusCode)
	}

	// We could decode the response body here if needed, but for basic tracking,
	// we usually just care about the status code.

	return nil
}
