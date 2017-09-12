package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CrossEngage/CaOps/internal/agent"
	"github.com/gorilla/mux"
)

// Server encapsulates all HTTP API
type Server struct {
	stopChan chan os.Signal
	server   *http.Server
	router   *mux.Router
	agent    *agent.Agent
}

// NewServer ...
func NewServer(bindTo string, agent *agent.Agent) *Server {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	router := mux.NewRouter()
	server := &Server{
		stopChan: stopChan,
		server:   &http.Server{Addr: bindTo, Handler: router},
		router:   router,
		agent:    agent,
	}

	router.Methods("GET").Path("/snapshot/{keyspaceGlob}/{table}").HandlerFunc(server.snapshotHandler)

	return server
}

func (s *Server) waitForShutdown() {
	<-s.stopChan
	log.Print("Shutting down HTTP server...")
	// shut down gracefully, but wait no longer than 5 seconds before halting
	// TODO make this configurable - maybe increase it for when there are uploads happening
	ctx, cancelFun := context.WithTimeout(context.Background(), 5*time.Second)
	cancelFun()
	s.server.Shutdown(ctx)
	log.Print("HTTP Server gracefully stopped")
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	go s.waitForShutdown()
	return s.server.ListenAndServe()
}
