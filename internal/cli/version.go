package cli

import (
	"fmt"
	"io"

	"github.com/nazar256/datadog-cli/internal/output"
	"github.com/spf13/cobra"
)

type versionView struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func newVersionCmd(opts *GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the CLI version",
		Args:    cobra.NoArgs,
		GroupID: "utility",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := runtimeConfigForOffline(opts)
			if err != nil {
				return err
			}
			view := versionView{Version: opts.BuildInfo.Version, Commit: opts.BuildInfo.Commit, Date: opts.BuildInfo.Date}
			return output.Write(cmd.OutOrStdout(), cfg.Output, view, func(w io.Writer) error {
				_, err := fmt.Fprintf(w, "ddog version %s (commit: %s, date: %s)\n", opts.BuildInfo.Version, opts.BuildInfo.Commit, opts.BuildInfo.Date)
				return err
			})
		},
	}
	return cmd
}
