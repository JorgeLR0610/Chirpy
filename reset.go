package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You do not have the right permissions")
		return
	}
	if err := cfg.DB.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("There was an error: %v", err))
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and Users table truncated"))
}
