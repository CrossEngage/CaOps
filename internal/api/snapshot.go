package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) snapshotHandler(w http.ResponseWriter, r *http.Request) {
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
