package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CrossEngage/CaOps/internal/agent"
	"github.com/gorilla/mux"
)

// HTTPService encapsulates all HTTP API
type HTTPService struct {
	stopChan chan os.Signal
	server   *http.Server
	router   *mux.Router
	agent    *agent.Agent
}

// NewHTTPService ...
func NewHTTPService(bindTo string, agent *agent.Agent) *HTTPService {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	router := mux.NewRouter()
	server := &HTTPService{
		stopChan: stopChan,
		server:   &http.Server{Addr: bindTo, Handler: router},
		router:   router,
		agent:    agent,
	}

	router.Methods("GET").Path("/snapshot/{keyspaceGlob}/{table}").HandlerFunc(server.snapshotHandler)

	return server
}

func (s *HTTPService) waitForShutdown() {
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
func (s *HTTPService) ListenAndServe() error {
	go s.waitForShutdown()
	return s.server.ListenAndServe()
}

func (s *HTTPService) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	defer r.Body.Close()
	// TODO trigger a IC snapshot, that will also check the cluster status
	log.Printf("Snapshot of %s.%s requested", vars["keyspaceGlob"], vars["table"])
	// TODO serfCli.UserEvent("Snapshot", []byte(payload), true)
	err := s.agent.DoSnapshot(vars["keyspaceGlob"], vars["table"])
	if err != nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Snapshot of %s.%s requested", vars["keyspaceGlob"], vars["table"])
}
