package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ensamblatec/CachSentinel/internal/core/service"
)

type ProxyHandler struct {
	cacheSvc *service.CacheService[any]
}

func NewProxyHandler(svc *service.CacheService[any]) *ProxyHandler {
	return &ProxyHandler{cacheSvc: svc}
}

func (handler *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "missing_resource_key", http.StatusBadRequest)
		return
	}

	data, err := handler.cacheSvc.Execute(r.Context(), key)
	if err != nil {
		http.Error(w, "upstream_error: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache-Sentinel", "HIT-OR-PREDICTIVE")
	json.NewEncoder(w).Encode(data)
}
