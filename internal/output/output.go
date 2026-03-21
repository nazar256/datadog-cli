package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode"
)

type Format string

const (
	Text Format = "text"
	JSON Format = "json"
)

func ParseFormat(raw string) (Format, error) {
	switch Format(strings.ToLower(strings.TrimSpace(raw))) {
	case "", Text:
		return Text, nil
	case JSON:
		return JSON, nil
	default:
		return "", fmt.Errorf("unsupported output format %q", raw)
	}
}

func Write(w io.Writer, format Format, value any, textRenderer func(io.Writer) error) error {
	if format == JSON {
		encoder := json.NewEncoder(w)
		return encoder.Encode(value)
	}
	return textRenderer(w)
}

func Table(w io.Writer, headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, strings.Join(sanitizeSlice(headers), "\t")); err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := fmt.Fprintln(tw, strings.Join(sanitizeSlice(row), "\t")); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func KeyValue(w io.Writer, pairs [][2]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, pair := range pairs {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", sanitizeText(pair[0]), sanitizeText(pair[1])); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func sanitizeSlice(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, sanitizeText(item))
	}
	return result
}

func sanitizeText(value string) string {
	var builder strings.Builder
	for _, r := range value {
		switch {
		case r == '\n':
			builder.WriteString(`\n`)
		case r == '\r':
			builder.WriteString(`\r`)
		case r == '\t':
			builder.WriteString(`\t`)
		case unicode.IsControl(r):
			builder.WriteString(`\x`)
			builder.WriteString(strings.ToUpper(strconv.FormatInt(int64(r), 16)))
		default:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
