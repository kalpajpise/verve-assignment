package router

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/kalpaj/verve/internal/config"
	"github.com/kalpaj/verve/pkg/db/redis"
)

type Router struct {
	cfg    *config.Config
	router *mux.Router
	mutex  *sync.RWMutex
	redis  *redis.Redis
}

func New(cfg *config.Config, redis *redis.Redis) (*Router, error) {

	if cfg == nil {
		return nil, fmt.Errorf("configuration is empty")
	}

	e := &Router{
		cfg:   cfg,
		mutex: &sync.RWMutex{},
		redis: redis,
	}
	e.init()

	return e, nil
}

func (e *Router) init() {
	e.router = mux.NewRouter()

	// Root path
	e.router.HandleFunc("/", e.root)
	e.router.HandleFunc("/v1/verve/accept", e.AcceptID).Methods("GET")
}

func (e *Router) root(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (e *Router) GetHandler() http.Handler {
	return e.router
}
