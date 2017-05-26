package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func httpServer() {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	router := mux.NewRouter()
	router.Methods("GET").Path("/snapshot/{keyspace}/{table}").HandlerFunc(snapshotHandler)

	srv := &http.Server{Addr: ":8081", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("Error while starting HTTP server", err)
		}
	}()

	<-stopChan // wait for SIGINT
	log.Info("Shutting down HTTP server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancelFun := context.WithTimeout(context.Background(), 5*time.Second)
	cancelFun()
	srv.Shutdown(ctx)

	log.Info("HTTP Server gracefully stopped")
}

func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := checkClusterStatus(); err != nil {
		http.Error(w, err.Error(), 500)
	}
	payload := fmt.Sprintf("%s.%s", vars["keyspace"], vars["table"])
	serfCli.UserEvent("Snapshot", []byte(payload), true)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Starting snapshot of %s.%s!\n\n", vars["keyspace"], vars["table"])
}
