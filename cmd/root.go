package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tcpgoon",
	Short: "tcpgoon tests concurrent connections towards a server listening on a TCP port",
	Long:  ``,
}

func Execute(buildstamp string, githash string) {
	releaseInfo.buildstamp = buildstamp
	releaseInfo.githash = githash

	AddCommands()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func AddCommands() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
}
