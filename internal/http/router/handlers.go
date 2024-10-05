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

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err != nil
}
