package api

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/JorgeLR0610/Chirpy/internal/service"
)


type AdminHandler struct {
	service *service.AdminService
	platform string
	fileServerHits	*atomic.Int32
}

func NewAdminHandler(svc *service.AdminService, platform string, hits *atomic.Int32) *AdminHandler {
	return &AdminHandler{service: svc, platform: platform, fileServerHits: hits}
}

func (h *AdminHandler) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>

	</html>
		`, h.fileServerHits.Load())))
}

func (h *AdminHandler) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (h *AdminHandler) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if h.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You do not have the right permissions")
		return		
	}

	if err := h.service.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not reset")
		return
	}

	h.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and Users table truncated"))	
}