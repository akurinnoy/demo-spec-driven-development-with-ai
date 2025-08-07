package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand" // Note: Go 1.20+ auto-seeds this package.
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// FR001: Word Generation: Two hard-coded string slices.
var adjectives = []string{
	"quick", "lazy", "sleepy", "noisy", "hungry", "brave", "bright", "calm", "eager", "fancy",
	"gentle", "happy", "jolly", "kind", "lively", "merry", "nice", "proud", "silly", "witty",
	"clever", "dizzy", "grumpy", "lucky", "mighty", "plucky", "shiny", "tiny", "zany", "breezy",
	"bubbly", "cheery", "comfy", "cozy", "crispy", "curly", "fluffy", "fuzzy", "gloomy", "groovy",
	"hazy", "icy", "jazzy", "jumpy", "quirky", "snappy", "sunny", "tasty", "vibey", "zippy",
}

var nouns = []string{
	"fox", "dog", "cat", "mouse", "bird", "wolf", "lion", "tiger", "bear", "frog",
	"fish", "shark", "whale", "squid", "crab", "snake", "lizard", "gecko", "newt", "otter",
	"seal", "duck", "goose", "swan", "crow", "raven", "owl", "hawk", "eagle", "pigeon",
	"robin", "finch", "sparrow", "horse", "pony", "donkey", "mule", "cow", "bull", "pig",
	"sheep", "goat", "llama", "alpaca", "camel", "koala", "panda", "sloth", "lemur", "hippo",
}

// Data Model for urls.json
type URLRecord struct {
	ShortCode  string `json:"short_code"`
	LongURL    string `json:"long_url"`
	CreatedAt  string `json:"created_at"`
	UsageCount int    `json:"usage_count"`
}

// In-memory cache of URL records, loaded from urls.json
var urls []URLRecord

// FR003: Concurrency Control: Mutex for all urls.json read/write operations.
var urlMutex = &sync.Mutex{}

var jsonFilePath = "urls.json"

// init function to load data on startup.
func init() {
	urlMutex.Lock()
	defer urlMutex.Unlock()

	// FR002: Data Persistence: Check for urls.json and create if not exists.
	if _, err := os.Stat(jsonFilePath); os.IsNotExist(err) {
		log.Printf("'%s' not found, creating it with default empty array.", jsonFilePath)
		if err := os.WriteFile(jsonFilePath, []byte("[]"), 0644); err != nil {
			log.Fatalf("Failed to create %s: %v", jsonFilePath, err)
		}
	}

	// Read the entire file.
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", jsonFilePath, err)
	}

	// Unmarshal the JSON data into the urls slice.
	if err := json.Unmarshal(data, &urls); err != nil {
		log.Fatalf("Failed to unmarshal JSON from %s: %v", jsonFilePath, err)
	}

	log.Printf("Loaded %d URL records from %s", len(urls), jsonFilePath)
}

// saveURLs writes the current state of the urls slice to urls.json.
// This function assumes the caller has already locked the mutex.
func saveURLs() error {
	// Marshal with indentation for readability.
	data, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal urls: %w", err)
	}
	if err := os.WriteFile(jsonFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %w", jsonFilePath, err)
	}
	return nil
}

// writeJSONError is a helper for sending consistent JSON error responses.
func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// createURLHandler handles POST /api/urls.
// FR004: API Endpoint POST /api/urls.
func createURLHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.URL == "" {
		writeJSONError(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Validate that the URL is in a reasonable format.
	if _, err := url.ParseRequestURI(reqBody.URL); err != nil {
		writeJSONError(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	urlMutex.Lock()
	defer urlMutex.Unlock()

	// Generate a unique short code that doesn't already exist.
	var shortCode string
	for {
		adj := adjectives[rand.Intn(len(adjectives))]
		noun := nouns[rand.Intn(len(nouns))]
		code := fmt.Sprintf("%s-%s", adj, noun)

		isUnique := true
		for _, record := range urls {
			if record.ShortCode == code {
				isUnique = false
				break
			}
		}
		if isUnique {
			shortCode = code
			break
		}
	}

	record := URLRecord{
		ShortCode:  shortCode,
		LongURL:    reqBody.URL,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		UsageCount: 0,
	}

	urls = append(urls, record)

	if err := saveURLs(); err != nil {
		log.Printf("Error saving URLs: %v", err)
		// FR008: Error Handling for file I/O.
		writeJSONError(w, "Failed to save URL record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"short_code": shortCode})
}

// getURLsHandler handles GET /api/urls.
// FR005: API Endpoint GET /api/urls.
func getURLsHandler(w http.ResponseWriter, r *http.Request) {
	urlMutex.Lock()
	defer urlMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	// Return a copy to avoid race conditions if the caller modifies the slice.
	urlsCopy := make([]URLRecord, len(urls))
	copy(urlsCopy, urls)

	if err := json.NewEncoder(w).Encode(urlsCopy); err != nil {
		log.Printf("Error encoding URLs: %v", err)
		// FR008: Error Handling.
		writeJSONError(w, "Failed to encode URL list", http.StatusInternalServerError)
	}
}

// rootHandler dispatches requests to the correct handler based on the URL path.
func rootHandler(staticFileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Route API calls.
		if strings.HasPrefix(r.URL.Path, "/api/urls") {
			switch r.Method {
			case http.MethodGet:
				getURLsHandler(w, r)
			case http.MethodPost:
				createURLHandler(w, r)
			default:
				writeJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// FR006: Redirect Endpoint GET /{short_code}.
		// Any path that is not an API call is a potential short code.
		shortCode := strings.TrimPrefix(r.URL.Path, "/")
		if shortCode != "" {
			urlMutex.Lock()
			var targetURL string
			found := false
			for i := range urls {
				if urls[i].ShortCode == shortCode {
					urls[i].UsageCount++
					targetURL = urls[i].LongURL
					found = true
					// Persist the change in usage_count.
					if err := saveURLs(); err != nil {
						urlMutex.Unlock()
						log.Printf("Error saving URLs on redirect: %v", err)
						writeJSONError(w, "Failed to update URL data", http.StatusInternalServerError)
						return
					}
					break
				}
			}
			urlMutex.Unlock()

			if found {
				http.Redirect(w, r, targetURL, http.StatusFound) // 302 Found
				return
			}
		}

		// FR007: Static File Server.
		// If not an API call and not a valid short code, fall back to serving static files.
		// This will correctly return a 404 if the file doesn't exist.
		staticFileServer.ServeHTTP(w, r)
	}
}

func main() {
	// FR007: The static file server serves assets from the frontend build directory.
	staticFileServer := http.FileServer(http.Dir("frontend/build/"))

	http.HandleFunc("/", rootHandler(staticFileServer))

	port := "8080"
	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
