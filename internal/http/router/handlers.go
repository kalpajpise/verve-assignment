package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/kalpaj/verve/pkg/constant"
)

func (e *Router) AcceptID(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	endpoint := r.URL.Query().Get("endpoint")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed"))
		return
	}

	log.Printf("Enpoint daata %s", endpoint)

	// Store the id in distributed redis set
	e.redis.SetAdd(constant.UniqueIDSet, id)

	if endpoint != "" {

		// Check if endpoint is a valid url
		if !isValidURL(endpoint) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed"))
			return
		}

		// Call the endpoint
		count, _ := e.redis.SetLength(constant.UniqueIDSet)
		fireUniqueIDCount(endpoint, count)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func fireUniqueIDCount(endpoint string, count int64) {
	payload := map[string]any{"count": count}
	body, _ := json.Marshal(payload)

	// Add the count as query param as well
	endpoint = fmt.Sprintf("%s?count=%d", strings.TrimPrefix(endpoint, "/"), count)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error sending POST request to %s: %v", endpoint, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("POST [%s] Status [%d]", endpoint, resp.StatusCode)
}

func isValidURL(urlStr string) bool {
	// Check if the URL string is empty
	if urlStr == "" {
		return false
	}

	// Parse the URL
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	// Check if scheme is present and is either http or https
	if u.Scheme == "" || !(u.Scheme == "http" || u.Scheme == "https") {
		return false
	}

	// Check if host is present
	if u.Host == "" {
		return false
	}

	// Check if the host has at least one dot and no underscores
	// This helps validate that it's a proper domain name
	hostParts := strings.Split(u.Host, ":")
	hostname := hostParts[0]
	if !strings.Contains(hostname, ".") || strings.Contains(hostname, "_") {
		return false
	}

	// Optionally, check if port is valid if present
	if len(hostParts) > 1 {
		port := hostParts[1]
		for _, ch := range port {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}

	return true
}
