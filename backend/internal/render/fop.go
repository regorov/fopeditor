package render

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FOPRenderer calls an external HTTP service that wraps Apache FOP.
type FOPRenderer struct {
	endpoint string
	client   *http.Client
}

// NewFOPRenderer builds a renderer targeting the provided endpoint.
func NewFOPRenderer(endpoint string) *FOPRenderer {
	return &FOPRenderer{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type fopRequest struct {
	XSL string `json:"xsl"`
	XML string `json:"xml"`
}

// Render forwards the payload to the FOP sidecar service and streams the PDF back.
func (r *FOPRenderer) Render(ctx context.Context, xsl, xml string) ([]byte, error) {
	payload, err := json.Marshal(fopRequest{XSL: xsl, XML: xml})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/pdf")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call fop sidecar: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read fop response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fop error: %s", bytes.TrimSpace(body))
	}

	return body, nil
}
