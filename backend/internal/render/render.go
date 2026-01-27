package render

import "context"

// Renderer defines the contract for PDF rendering implementations.
type Renderer interface {
	Render(ctx context.Context, xsl, xml string) ([]byte, error)
}
