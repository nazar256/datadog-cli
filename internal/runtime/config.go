package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/nazar256/datadog-cli/internal/output"
)

const (
	DefaultSite    = "datadoghq.com"
	DefaultTimeout = 30 * time.Second
)

type FlagValues struct {
	Site      string
	EnvFile   string
	NoEnvFile bool
	Timeout   time.Duration
	Output    string
}

type Config struct {
	Site        string
	Timeout     time.Duration
	Output      output.Format
	APIKey      string
	AppKey      string
	EnvFileUsed string
	Version     string
}

func ResolveConfig(flags FlagValues) (Config, error) {
	dotenv, envFileUsed, err := readDotEnv(flags)
	if err != nil {
		return Config{}, err
	}

	lookup := func(key string) string {
		if value, ok := os.LookupEnv(key); ok {
			return strings.TrimSpace(value)
		}
		return strings.TrimSpace(dotenv[key])
	}

	site, err := normalizeSite(firstNonEmpty(flags.Site, lookup("DATADOG_SITE"), DefaultSite))
	if err != nil {
		return Config{}, err
	}

	format, err := output.ParseFormat(flags.Output)
	if err != nil {
		return Config{}, err
	}

	timeout := flags.Timeout
	if timeout < 0 {
		return Config{}, fmt.Errorf("timeout must be greater than or equal to 0")
	}
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return Config{
		Site:        site,
		Timeout:     timeout,
		Output:      format,
		APIKey:      lookup("DATADOG_API_KEY"),
		AppKey:      lookup("DATADOG_APP_KEY"),
		EnvFileUsed: envFileUsed,
	}, nil
}

func (c Config) HasAuth() bool {
	return c.APIKey != "" && c.AppKey != ""
}

func (c Config) RequireAuth() error {
	if c.APIKey == "" {
		return fmt.Errorf("missing DATADOG_API_KEY")
	}
	if c.AppKey == "" {
		return fmt.Errorf("missing DATADOG_APP_KEY")
	}
	return nil
}

func readDotEnv(flags FlagValues) (map[string]string, string, error) {
	if flags.NoEnvFile {
		return map[string]string{}, "", nil
	}

	path := strings.TrimSpace(flags.EnvFile)
	if path == "" {
		path = ".env"
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve env file path: %w", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) && flags.EnvFile == "" {
			return map[string]string{}, "", nil
		}
		return nil, "", fmt.Errorf("read env file: %w", err)
	}

	values, err := godotenv.Read(absPath)
	if err != nil {
		return nil, "", fmt.Errorf("parse env file: %w", err)
	}
	return values, absPath, nil
}

func normalizeSite(raw string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	aliases := map[string]string{
		"us1":     "datadoghq.com",
		"us3":     "us3.datadoghq.com",
		"us5":     "us5.datadoghq.com",
		"eu":      "datadoghq.eu",
		"ap1":     "ap1.datadoghq.com",
		"ap2":     "ap2.datadoghq.com",
		"us1-fed": "ddog-gov.com",
	}
	if value == "" {
		return "", fmt.Errorf("site cannot be empty")
	}
	if normalized, ok := aliases[value]; ok {
		return normalized, nil
	}
	allowed := map[string]struct{}{
		"datadoghq.com":     {},
		"us3.datadoghq.com": {},
		"us5.datadoghq.com": {},
		"datadoghq.eu":      {},
		"ap1.datadoghq.com": {},
		"ap2.datadoghq.com": {},
		"ddog-gov.com":      {},
	}
	if _, ok := allowed[value]; ok {
		return value, nil
	}
	return "", fmt.Errorf("unsupported Datadog site %q", raw)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
