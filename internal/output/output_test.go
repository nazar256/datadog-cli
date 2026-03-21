package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestSanitizeText(t *testing.T) {
	got := sanitizeText("hello\nworld\x1b[31m\tend")
	if strings.Contains(got, "\x1b") || strings.Contains(got, "\nworld\n") {
		t.Fatalf("unexpected unsafe output: %q", got)
	}
	if !strings.Contains(got, `\n`) || !strings.Contains(got, `\t`) || !strings.Contains(got, `\x1B`) {
		t.Fatalf("expected escaped control characters, got %q", got)
	}
}

func TestTableSanitizesCells(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Table(buf, []string{"NAME"}, [][]string{{"line1\nline2\x1b[31m"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "\x1b") || strings.Contains(out, "line1\nline2") {
		t.Fatalf("unexpected unsafe table output: %q", out)
	}
}
