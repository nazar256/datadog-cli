package runtime

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nazar256/datadog-cli/internal/timeutil"
)

func TestResolveConfig(t *testing.T) {
	// Create a temporary .env file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envFile, []byte("DATADOG_SITE=eu\nDATADOG_API_KEY=env_api\nDATADOG_APP_KEY=env_app\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}

	tests := []struct {
		name       string
		flags      FlagValues
		envVars    map[string]string
		wantSite   string
		wantAPIKey string
		wantAppKey string
		wantErr    bool
	}{
		{
			name: "defaults",
			flags: FlagValues{
				NoEnvFile: true,
			},
			wantSite: DefaultSite,
		},
		{
			name: "flag overrides env var",
			flags: FlagValues{
				Site:      "us3",
				NoEnvFile: true,
			},
			envVars: map[string]string{
				"DATADOG_SITE": "eu",
			},
			wantSite: "us3.datadoghq.com",
		},
		{
			name: "env var overrides .env file",
			flags: FlagValues{
				EnvFile: envFile,
			},
			envVars: map[string]string{
				"DATADOG_SITE":    "us5",
				"DATADOG_API_KEY": "var_api",
			},
			wantSite:   "us5.datadoghq.com",
			wantAPIKey: "var_api",
			wantAppKey: "env_app",
		},
		{
			name: "empty env masks dotenv values",
			flags: FlagValues{
				EnvFile: envFile,
			},
			envVars: map[string]string{
				"DATADOG_SITE":    "",
				"DATADOG_API_KEY": "",
				"DATADOG_APP_KEY": "",
			},
			wantSite: DefaultSite,
		},
		{
			name: "no env file flag ignores .env",
			flags: FlagValues{
				EnvFile:   envFile,
				NoEnvFile: true,
			},
			wantSite: DefaultSite,
		},
		{
			name: "invalid site",
			flags: FlagValues{
				Site:      "invalid",
				NoEnvFile: true,
			},
			wantErr: true,
		},
		{
			name: "reject unknown dotted host",
			flags: FlagValues{
				Site:      "evil.example.com",
				NoEnvFile: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := ResolveConfig(tt.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if cfg.Site != tt.wantSite {
				t.Errorf("ResolveConfig() Site = %v, want %v", cfg.Site, tt.wantSite)
			}
			if tt.wantAPIKey != "" && cfg.APIKey != tt.wantAPIKey {
				t.Errorf("ResolveConfig() APIKey = %v, want %v", cfg.APIKey, tt.wantAPIKey)
			}
			if tt.wantAPIKey == "" && cfg.APIKey != "" {
				t.Errorf("ResolveConfig() APIKey = %v, want empty", cfg.APIKey)
			}
			if tt.wantAppKey != "" && cfg.AppKey != tt.wantAppKey {
				t.Errorf("ResolveConfig() AppKey = %v, want %v", cfg.AppKey, tt.wantAppKey)
			}
			if tt.wantAppKey == "" && cfg.AppKey != "" {
				t.Errorf("ResolveConfig() AppKey = %v, want empty", cfg.AppKey)
			}
		})
	}
}

func TestNormalizeSite(t *testing.T) {
	tests := []struct {
		raw     string
		want    string
		wantErr bool
	}{
		{"us1", "datadoghq.com", false},
		{"us3", "us3.datadoghq.com", false},
		{"eu", "datadoghq.eu", false},
		{"us5.datadoghq.com", "us5.datadoghq.com", false},
		{"invalid", "", true},
		{"custom.datadoghq.com", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			got, err := normalizeSite(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeSite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeSite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRangeRejectsNonPositiveLast(t *testing.T) {
	_, err := timeutil.ParseRange("-15m", "", "", func() time.Time { return time.Unix(0, 0) })
	if err == nil {
		t.Fatal("expected error for negative --last")
	}
	_, err = timeutil.ParseRange("0s", "", "", func() time.Time { return time.Unix(0, 0) })
	if err == nil {
		t.Fatal("expected error for zero --last")
	}
}

func TestResolveConfigRejectsNegativeTimeout(t *testing.T) {
	_, err := ResolveConfig(FlagValues{NoEnvFile: true, Timeout: -1 * time.Second})
	if err == nil {
		t.Fatal("expected negative timeout to fail")
	}
}
