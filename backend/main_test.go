package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// setupTest creates a temporary urls.json for testing and returns a teardown function.
// It resets the in-memory store and file path for each test, ensuring isolation.
func setupTest(t *testing.T) func() {
	// Create a temporary file for urls.json
	tmpfile, err := os.CreateTemp("", "urls.*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write initial empty array to the temp file
	if _, err := tmpfile.Write([]byte("[]")); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Lock the mutex to safely change shared state
	urlMutex.Lock()

	// Override the global jsonFilePath for the duration of the test
	originalPath := jsonFilePath
	jsonFilePath = tmpfile.Name()

	// Reset the in-memory store
	urls = []URLRecord{}

	// Manually load from the new temp file, mimicking the behavior of init()
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		urlMutex.Unlock()
		t.Fatalf("Failed to read temp json file: %v", err)
	}
	if err := json.Unmarshal(data, &urls); err != nil {
		urlMutex.Unlock()
		t.Fatalf("Failed to unmarshal temp json data: %v", err)
	}

	urlMutex.Unlock()

	// Return a teardown function to be called at the end of the test
	return func() {
		urlMutex.Lock()
		defer urlMutex.Unlock()
		os.Remove(tmpfile.Name())
		jsonFilePath = originalPath // Restore original path
		urls = []URLRecord{}         // Clear in-memory store
	}
}

func TestCreateURL_Success(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	handler := rootHandler(http.NotFoundHandler())

	body := `{"url": "https://example.com/a-very-long-url"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var respBody map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Could not unmarshal response body: %v", err)
	}
	if _, ok := respBody["short_code"]; !ok {
		t.Errorf("response body does not contain short_code")
	}

	urlMutex.Lock()
	defer urlMutex.Unlock()
	if len(urls) != 1 {
		t.Errorf("expected 1 URL in memory, got %d", len(urls))
	}
	if urls[0].LongURL != "https://example.com/a-very-long-url" {
		t.Errorf("wrong long_url stored in memory")
	}
}

func TestCreateURL_InvalidURL(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	handler := rootHandler(http.NotFoundHandler())

	body := `{"url": "not-a-valid-url"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var respBody map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Could not unmarshal response body: %v", err)
	}
	expectedError := "Invalid URL format"
	if respBody["error"] != expectedError {
		t.Errorf("unexpected error message: got '%v' want '%v'", respBody["error"], expectedError)
	}
}

func TestRedirect_Success(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Manually add a URL to the store for the test
	urlMutex.Lock()
	testRecord := URLRecord{
		ShortCode:  "test-code",
		LongURL:    "https://example.com/redirect-target",
		UsageCount: 0,
	}
	urls = append(urls, testRecord)
	if err := saveURLs(); err != nil {
		urlMutex.Unlock()
		t.Fatalf("Failed to save test URL: %v", err)
	}
	urlMutex.Unlock()

	handler := rootHandler(http.NotFoundHandler())

	req := httptest.NewRequest(http.MethodGet, "/test-code", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	location := rr.Header().Get("Location")
	expectedLocation := "https://example.com/redirect-target"
	if location != expectedLocation {
		t.Errorf("handler returned wrong redirect location: got %v want %v", location, expectedLocation)
	}

	// Check usage count increment
	urlMutex.Lock()
	defer urlMutex.Unlock()
	if len(urls) != 1 || urls[0].UsageCount != 1 {
		t.Errorf("usage count was not incremented: got %d", urls[0].UsageCount)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	handler := rootHandler(http.NotFoundHandler())

	req := httptest.NewRequest(http.MethodGet, "/non-existent-code", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// The request should fall through to the static file server, which returns 404
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestGetURLs(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Manually add a URL to the store for the test
	urlMutex.Lock()
	testRecord := URLRecord{
		ShortCode:  "test-code",
		LongURL:    "https://example.com/get-urls-test",
		UsageCount: 5,
	}
	urls = append(urls, testRecord)
	if err := saveURLs(); err != nil {
		urlMutex.Unlock()
		t.Fatalf("Failed to save test URL: %v", err)
	}
	urlMutex.Unlock()

	handler := rootHandler(http.NotFoundHandler())

	req := httptest.NewRequest(http.MethodGet, "/api/urls", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var respBody []URLRecord
	if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Could not unmarshal response body: %v", err)
	}

	if len(respBody) != 1 {
		t.Fatalf("expected 1 URL in response, got %d", len(respBody))
	}
	if respBody[0].ShortCode != "test-code" || respBody[0].LongURL != "https://example.com/get-urls-test" || respBody[0].UsageCount != 5 {
		t.Errorf("unexpected URL data in response: got %+v", respBody[0])
	}
}
