package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type renderRequest struct {
	XSL string `json:"xsl"`
	XML string `json:"xml"`
}

func main() {
	port := getenv("PORT", "8090")
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/render", handleRender)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("FOP sidecar listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleRender(w http.ResponseWriter, r *http.Request) {
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
	if req.XSL == "" || req.XML == "" {
		http.Error(w, "xsl and xml are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	pdf, err := renderWithFOP(ctx, req.XSL, req.XML)
	if err != nil {
		log.Printf("render failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=render.pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdf)
}

func renderWithFOP(ctx context.Context, xsl, xml string) ([]byte, error) {
	dir, err := os.MkdirTemp("", "fop-work-")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	xslPath := filepath.Join(dir, "template.xsl")
	xmlPath := filepath.Join(dir, "data.xml")
	pdfPath := filepath.Join(dir, "output.pdf")

	if err := os.WriteFile(xslPath, []byte(xsl), 0o600); err != nil {
		return nil, fmt.Errorf("write xsl: %w", err)
	}
	if err := os.WriteFile(xmlPath, []byte(xml), 0o600); err != nil {
		return nil, fmt.Errorf("write xml: %w", err)
	}

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "fop", "-xml", xmlPath, "-xsl", xslPath, "-pdf", pdfPath)
	cmd.Stdout = &stderr
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("fop command failed: %w: %s", err, stderr.String())
	}

	output, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("read pdf: %w", err)
	}

	return output, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
