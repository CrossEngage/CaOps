package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (caops *CaOps) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	defer r.Body.Close()
	log.Printf("Snapshot of %s.%s requested", vars["keyspaceGlob"], vars["table"])
	payload := fmt.Sprintf("%s:%s", vars["keyspaceGlob"], vars["table"])
	err := caops.gossiper.SendEvent("SnapshotGlob", payload)
	if err != nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Snapshot of %s.%s requested", vars["keyspaceGlob"], vars["table"])
}
