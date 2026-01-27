package render

import (
	"bytes"
	"context"
	"fmt"
)

// StubRenderer implements Renderer by returning a placeholder PDF.
type StubRenderer struct{}

// NewStubRenderer creates a renderer suitable for local development
// when Apache FOP is not yet wired up.
func NewStubRenderer() *StubRenderer {
	return &StubRenderer{}
}

// Render generates a static PDF that references the provided payload sizes.
func (s *StubRenderer) Render(ctx context.Context, xsl, xml string) ([]byte, error) {
	_ = ctx
	content := fmt.Sprintf("XSL bytes: %d\nXML bytes: %d", len(xsl), len(xml))
	return buildSimplePDF(content), nil
}

func buildSimplePDF(content string) []byte {
	text := pdfEscape(content)
	stream := fmt.Sprintf("BT /F1 14 Tf 72 720 Td (%s) Tj ET", text)

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")

	offsets := []int{0}
	writeObj := func(obj string) {
		offsets = append(offsets, buf.Len())
		buf.WriteString(obj)
	}

	writeObj("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	writeObj("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	writeObj("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n")
	writeObj(fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(stream), stream))
	writeObj("5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	xrefOffset := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n", len(offsets))
	buf.WriteString("0000000000 65535 f \n")
	for _, off := range offsets[1:] {
		fmt.Fprintf(&buf, "%010d 00000 n \n", off)
	}
	buf.WriteString("trailer\n")
	buf.WriteString("<< /Size 6 /Root 1 0 R >>\n")
	fmt.Fprintf(&buf, "startxref\n%d\n%%%%EOF", xrefOffset)

	return buf.Bytes()
}

func pdfEscape(s string) string {
	replacements := map[rune]string{
		'\\': "\\\\",
		'(': "\\(",
		')': "\\)",
	}
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if repl, ok := replacements[r]; ok {
			for _, rr := range repl {
				out = append(out, rr)
			}
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
