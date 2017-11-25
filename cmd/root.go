package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type tcpgoonFlags struct {
	hostPtr              string
	portPtr              int
	numberConnectionsPtr int
	delayPtr             int
	connDialTimeoutPtr   int
	debugPtr             bool
	reportingIntervalPtr int
	assumeyesPtr         bool
}

var flags tcpgoonFlags

var rootCmd = &cobra.Command{
	Use:   "tcpgoon",
	Short: "tcpgoon tests concurrent connections towards a server listening on a TCP port",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.NoArgs(cmd, args)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := validateFlags(flags); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		run(flags)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flags.hostPtr, "target", "t", "", "Target host you want to open tcp connections against (Required)")
	rootCmd.Flags().IntVarP(&flags.portPtr, "port", "p", 0, "Port you want to open tcp connections against (Required)")
	rootCmd.Flags().IntVarP(&flags.numberConnectionsPtr, "connections", "c", 100, "Number of connections you want to open")
	rootCmd.Flags().IntVarP(&flags.delayPtr, "sleep", "s", 10, "Time you want to sleep between connections, in ms")
	rootCmd.Flags().IntVarP(&flags.connDialTimeoutPtr, "dial-timeout", "d", 5000, "Connection dialing timeout, in ms")
	rootCmd.Flags().BoolVarP(&flags.debugPtr, "verbose", "v", false, "Print debugging information to the standard error")
	rootCmd.Flags().IntVarP(&flags.reportingIntervalPtr, "interval", "i", 1, "Interval, in seconds, between stats updates")
	rootCmd.Flags().BoolVarP(&flags.assumeyesPtr, "assume-yes", "y", false, "Force execution without asking for confirmation")

}

func validateFlags(flags tcpgoonFlags) error {
	// Target host and port are mandatory to run the TCP check
	if flags.hostPtr == "" || flags.portPtr == 0 {
		return errors.New("Missing some required parameters")
	}
	return nil
}
