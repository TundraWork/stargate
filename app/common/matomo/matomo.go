package matomo

import (
	"bytes"
	"context"
	"github.com/cloudwego/hertz/pkg/common/json"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Event represents a single event to be tracked by Matomo.
type Event struct {
	ActionName string    // Required: The name of the action being tracked
	URL        string    // Required: The URL of the page/resource being accessed
	UserAgent  string    // Required: The user's browser user agent
	ClientIP   string    // Required: The user's IP address
	ClientTime time.Time // Required: The time of the event, in the client's timezone
}

type Client struct {
	matomoURL   string
	siteID      string
	authToken   string
	httpClient  *http.Client
	eventChan   chan Event
	stopChan    chan struct{}
	workerGroup sync.WaitGroup
	batchSize   int
}

var (
	clientInstance *Client
	once           sync.Once
)

// InitClient initializes the Matomo client.  It should be called once,
// typically during application startup.
func InitClient(matomoURL string, siteID string, authToken string, numWorkers int, batchSize int, eventBufferSize int) {
	once.Do(func() {
		clientInstance = &Client{
			matomoURL: matomoURL,
			siteID:    siteID,
			authToken: authToken,
			httpClient: &http.Client{
				Timeout: 5 * time.Second, // Set a reasonable timeout
			},
			eventChan: make(chan Event, eventBufferSize),
			stopChan:  make(chan struct{}),
			batchSize: batchSize,
		}

		if numWorkers <= 0 {
			numWorkers = 1 // Default to 1 worker if an invalid value is provided
		}

		clientInstance.workerGroup.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			go clientInstance.eventWorker(i)
		}
		hlog.Infof("[Matomo] Initialized client with %d workers, buffer size %d", numWorkers, eventBufferSize)
	})
}

// ReportEvent queues an event for reporting to Matomo.  It returns immediately.
func ReportEvent(ctx context.Context, event Event) {
	if clientInstance == nil {
		hlog.CtxErrorf(ctx, "[Matomo] ReportEvent called before InitMatomoClient")
		return
	}
	select {
	case clientInstance.eventChan <- event:
		// Event successfully queued
	default:
		hlog.CtxWarnf(ctx, "[Matomo] Event buffer full, dropping event: %v", event)
	}
}

// eventWorker is the worker goroutine that processes events from the channel.
func (c *Client) eventWorker(workerID int) {
	defer c.workerGroup.Done()
	hlog.Infof("[Matomo] Starting worker %d", workerID)

	// Use a ticker for batch processing
	ticker := time.NewTicker(1 * time.Second) // Send batches every second (or when full)
	defer ticker.Stop()

	var batch []Event
	for {
		select {
		case event := <-c.eventChan:
			batch = append(batch, event)
			if len(batch) >= c.batchSize {
				c.sendBatch(context.Background(), batch)
				batch = nil // Reset the batch
			}
		case <-ticker.C:
			if len(batch) > 0 {
				c.sendBatch(context.Background(), batch)
				batch = nil // Reset the batch
			}
		case <-c.stopChan:
			hlog.Infof("[Matomo] Stopping worker %d", workerID)
			// Send any remaining events before exiting
			if len(batch) > 0 {
				c.sendBatch(context.Background(), batch)
			}
			return
		}
	}
}

// sendEvent sends a single event to Matomo.
func (c *Client) sendBatch(ctx context.Context, events []Event) {
	if len(events) == 0 {
		return
	}

	requests := make([]string, len(events))
	for i, event := range events {
		params := url.Values{}
		params.Set("idsite", c.siteID)
		params.Set("rec", "1")
		params.Set("action_name", event.ActionName)
		params.Set("url", event.URL)
		params.Set("ua", event.UserAgent)
		params.Set("cip", event.ClientIP)
		params.Set("cdt", event.ClientTime.Format("2006-01-02 15:04:05"))
		params.Set("apiv", "1")
		params.Set("rand", strconv.FormatInt(time.Now().UnixNano(), 10))
		requests[i] = "?" + params.Encode()
	}

	payload := map[string]interface{}{
		"requests":   requests,
		"token_auth": c.authToken,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		hlog.CtxErrorf(ctx, "[Matomo] Error marshaling JSON: %v", err)
		return
	}

	resp, err := c.httpClient.Post(c.matomoURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		hlog.CtxErrorf(ctx, "[Matomo] Error sending batch: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		hlog.CtxErrorf(ctx, "[Matomo] Matomo returned non-OK status: %d", resp.StatusCode)
		return
	}

	hlog.CtxInfof(ctx, "[Matomo] Successfully sent batch of %d events", len(events))
}

// Shutdown gracefully shuts down the Matomo client, waiting for all pending events to be processed.
func Shutdown(ctx context.Context) {
	if clientInstance == nil {
		return // Nothing to shut down
	}

	hlog.CtxInfof(ctx, "[Matomo] Shutting down client...")
	close(clientInstance.stopChan) // Signal workers to stop

	done := make(chan struct{})
	go func() {
		clientInstance.workerGroup.Wait() // Wait for all workers to finish
		close(done)
	}()

	select {
	case <-done:
		hlog.CtxInfof(ctx, "[Matomo] Client shut down gracefully.")
	case <-time.After(10 * time.Second): // Timeout after a reasonable period
		hlog.CtxErrorf(ctx, "[Matomo] Client shutdown timed out.")
	}

	clientInstance = nil // Reset the instance
}
