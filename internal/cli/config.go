package cli

import (
	"fmt"
	"io"

	"github.com/nazar256/datadog-cli/internal/output"
	"github.com/nazar256/datadog-cli/internal/runtime"
	"github.com/spf13/cobra"
)

func newConfigCmd(opts *GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Inspect CLI configuration",
		GroupID: "utility",
	}

	cmd.AddCommand(newConfigDoctorCmd(opts))

	return cmd
}

func newConfigDoctorCmd(opts *GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "doctor",
		Short:   "Check configuration and authentication status",
		Args:    cobra.NoArgs,
		Long:    "Resolve configuration exactly as the CLI sees it and report non-secret settings plus whether authentication values are present.",
		Example: "ddog config doctor\n  ddog --site eu config doctor --output json\n  ddog --env-file .env config doctor",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := runtime.ResolveConfig(opts.FlagValues)
			if err != nil {
				return fmt.Errorf("failed to resolve config: %w", err)
			}

			doctor := configDoctorView{
				Site:       cfg.Site,
				Timeout:    cfg.Timeout.String(),
				Output:     string(cfg.Output),
				EnvFile:    emptyFallback(cfg.EnvFileUsed, "(none)"),
				APIKey:     presence(cfg.APIKey),
				AppKey:     presence(cfg.AppKey),
				AuthStatus: doctorStatus(cfg),
			}

			return output.Write(cmd.OutOrStdout(), cfg.Output, doctor, func(w io.Writer) error {
				_, err := fmt.Fprintln(w, "Configuration Doctor")
				if err != nil {
					return err
				}
				_, err = fmt.Fprintln(w, "--------------------")
				if err != nil {
					return err
				}
				if err := output.KeyValue(w, [][2]string{
					{"Site", doctor.Site},
					{"Timeout", doctor.Timeout},
					{"Output", doctor.Output},
					{"Env File", doctor.EnvFile},
					{"API Key", doctor.APIKey},
					{"App Key", doctor.AppKey},
					{"Status", doctor.AuthStatus},
				}); err != nil {
					return err
				}
				_, err = fmt.Fprintln(w, "\nSecrets are never printed. Use DATADOG_API_KEY and DATADOG_APP_KEY via env or .env.")
				return err
			})
		},
	}
	return cmd
}

type configDoctorView struct {
	Site       string `json:"site"`
	Timeout    string `json:"timeout"`
	Output     string `json:"output"`
	EnvFile    string `json:"env_file"`
	APIKey     string `json:"api_key"`
	AppKey     string `json:"app_key"`
	AuthStatus string `json:"auth_status"`
}

func runtimeConfigForOffline(opts *GlobalOptions) (runtime.Config, error) {
	format, err := output.ParseFormat(opts.Output)
	if err != nil {
		return runtime.Config{}, err
	}
	return runtime.Config{Output: format}, nil
}

func presence(value string) string {
	if value == "" {
		return "missing"
	}
	return "present"
}

func doctorStatus(cfg runtime.Config) string {
	if cfg.HasAuth() {
		return "ready"
	}
	return "missing_credentials"
}

func emptyFallback(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
