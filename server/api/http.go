package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
)

// Server encapsulates all HTTP API
type Server struct {
	log      *logging.Logger
	stopChan chan os.Signal
	server   *http.Server
	router   *mux.Router
}

// NewServer ...
func NewServer(log *logging.Logger, bindTo string) *Server {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	router := mux.NewRouter()
	server := &Server{
		log:      log,
		stopChan: stopChan,
		server:   &http.Server{Addr: bindTo, Handler: router},
		router:   router,
	}

	router.Methods("GET").Path("/snapshot/{keyspace}/{table}").HandlerFunc(server.snapshotHandler)

	return server
}

func (s *Server) waitForShutdown() {
	<-s.stopChan
	s.log.Info("Shutting down HTTP server...")
	// shut down gracefully, but wait no longer than 5 seconds before halting
	// TODO make this configurable - maybe increase it for when there are uploads happening
	ctx, cancelFun := context.WithTimeout(context.Background(), 5*time.Second)
	cancelFun()
	s.server.Shutdown(ctx)
	s.log.Info("HTTP Server gracefully stopped")
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	go s.waitForShutdown()
	return s.server.ListenAndServe()
}

func (s *Server) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// TODO
	// if err := checkClusterStatus(); err != nil {
	// 	http.Error(w, err.Error(), 500)
	// }
	// TODO
	// payload := fmt.Sprintf("%s.%s", vars["keyspace"], vars["table"])
	// serfCli.UserEvent("Snapshot", []byte(payload), true)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Starting snapshot of %s.%s!\n\n", vars["keyspace"], vars["table"])
}
