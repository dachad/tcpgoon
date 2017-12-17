package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show tcpgoon version",
	RunE: func(cmd *cobra.Command, args []string) error {
		printTcpgoonVersion()
		return nil
	},
}

func printTcpgoonVersion() {
	fmt.Println("Git Commit Hash: " + releaseInfo.githash)
	fmt.Println("UTC Build Time: " + releaseInfo.buildstamp)
}
