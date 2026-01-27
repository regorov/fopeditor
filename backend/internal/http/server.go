package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/regorov/fopeditor/backend/internal/render"
)

// Server wires HTTP handlers with the rendering backend.
type Server struct {
	renderer render.Renderer
}

// NewServer builds the HTTP API with the provided renderer.
func NewServer(renderer render.Renderer) *Server {
	return &Server{renderer: renderer}
}

// Handler exposes the configured HTTP mux for use in net/http servers.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/render", s.handleRender)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type renderRequest struct {
	XSL string `json:"xsl"`
	XML string `json:"xml"`
}

func (s *Server) handleRender(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req renderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}
	if err := validateRenderRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pdf, err := s.renderer.Render(r.Context(), req.XSL, req.XML)
	if err != nil {
		log.Printf("render error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=render.pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdf)
}

func validateRenderRequest(req renderRequest) error {
	if req.XSL == "" {
		return errors.New("xsl is required")
	}
	if req.XML == "" {
		return errors.New("xml is required")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
