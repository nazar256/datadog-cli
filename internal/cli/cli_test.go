package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmdHelp(t *testing.T) {
	cmd := NewRootCmd(BuildInfo{
		Version: "1.0.0",
		Commit:  "abcdef",
		Date:    "2023-10-27",
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check for command groups
	if !strings.Contains(output, "Core Commands:") {
		t.Errorf("expected help to contain 'Core Commands:', got:\n%s", output)
	}
	if !strings.Contains(output, "Utility Commands:") {
		t.Errorf("expected help to contain 'Utility Commands:', got:\n%s", output)
	}

	// Check for commands
	if !strings.Contains(output, "ddog is a Datadog CLI for humans, coding agents, and automation.") {
		t.Errorf("expected help to contain updated product description")
	}
	if !strings.Contains(output, "version") {
		t.Errorf("expected help to contain 'version' command")
	}
	if !strings.Contains(output, "monitor") {
		t.Errorf("expected help to contain 'monitor' command")
	}
	if !strings.Contains(output, "docs") {
		t.Errorf("expected help to contain 'docs' command")
	}
	if !strings.Contains(output, "config") {
		t.Errorf("expected help to contain 'config' command")
	}

	// Check for global flags
	if !strings.Contains(output, "--site") {
		t.Errorf("expected help to contain '--site' flag")
	}
	if !strings.Contains(output, "--env-file") {
		t.Errorf("expected help to contain '--env-file' flag")
	}
}

func TestVersionCmd(t *testing.T) {
	cmd := NewRootCmd(BuildInfo{
		Version: "1.0.0",
		Commit:  "abcdef",
		Date:    "2023-10-27",
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	expected := "ddog version 1.0.0 (commit: abcdef, date: 2023-10-27)\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestVersionCmdJSON(t *testing.T) {
	cmd := NewRootCmd(BuildInfo{
		Version: "1.0.0",
		Commit:  "abcdef",
		Date:    "2023-10-27",
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version", "--output", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"version":"1.0.0"`) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestConfigDoctorCmd(t *testing.T) {
	t.Setenv("DATADOG_API_KEY", "super-secret-api")
	t.Setenv("DATADOG_APP_KEY", "super-secret-app")

	cmd := NewRootCmd(BuildInfo{})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"config", "doctor", "--no-env-file"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check that it reports non-secret config
	if !strings.Contains(output, "Site      datadoghq.com") {
		t.Errorf("expected output to contain default site, got:\n%s", output)
	}

	// Check that it reports auth presence without leaking secrets
	if !strings.Contains(output, "API Key   present") {
		t.Errorf("expected output to report present API key, got:\n%s", output)
	}
	if !strings.Contains(output, "App Key   present") {
		t.Errorf("expected output to report present App key, got:\n%s", output)
	}

	if !strings.Contains(output, "Status    ready") {
		t.Errorf("expected output to report ready status, got:\n%s", output)
	}

	// Ensure it doesn't print actual secret values.
	if strings.Contains(output, "super-secret-api") || strings.Contains(output, "super-secret-app") {
		t.Errorf("output should not leak secret keys")
	}
}
