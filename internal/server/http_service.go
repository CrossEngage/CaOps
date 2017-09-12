package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// HTTPService encapsulates all HTTP API
type HTTPService struct {
	stopChan chan os.Signal
	server   *http.Server
	router   *mux.Router
	gossiper *Gossiper
}

// NewHTTPService ...
func NewHTTPService(bindTo string, gossiper *Gossiper) *HTTPService {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	router := mux.NewRouter()
	server := &HTTPService{
		stopChan: stopChan,
		server:   &http.Server{Addr: bindTo, Handler: router},
		router:   router,
		gossiper: gossiper,
	}

	router.Methods("GET").
		Path("/snapshot/{keyspaceGlob}/{table}").
		HandlerFunc(server.snapshotHandler)

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
	log.Printf("Snapshot of %s.%s requested", vars["keyspaceGlob"], vars["table"])
	payload := fmt.Sprintf("%s:%s", vars["keyspaceGlob"], vars["table"])
	err := s.gossiper.SendEvent("SnapshotGlob", payload)
	if err != nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Snapshot of %s.%s requested", vars["keyspaceGlob"], vars["table"])
}
