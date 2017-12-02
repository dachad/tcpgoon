package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dachad/tcpgoon/cmdutil"
	"github.com/dachad/tcpgoon/debugging"
	"github.com/spf13/cobra"
)

type tcpgoonParams struct {
	hostPtr              string
	portPtr              int
	numberConnectionsPtr int
	delayPtr             int
	connDialTimeoutPtr   int
	debugPtr             bool
	reportingIntervalPtr int
	assumeyesPtr         bool
}

var params tcpgoonParams

var rootCmd = &cobra.Command{
	Use:   "tcpgoon <host> <port>",
	Short: "tcpgoon tests concurrent connections towards a server listening on a TCP port",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := validateRequiredArgs(&params, args); err != nil {
			cmd.Println(cmd.UsageString())
			os.Exit(1)
		}
		enableDebugging(params)
		autorunValidation(params)
	},
	Run: func(cmd *cobra.Command, args []string) {
		run(params)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&params.numberConnectionsPtr, "connections", "c", 100, "Number of connections you want to open")
	rootCmd.Flags().IntVarP(&params.delayPtr, "sleep", "s", 10, "Time you want to sleep between connections, in ms")
	rootCmd.Flags().IntVarP(&params.connDialTimeoutPtr, "dial-timeout", "d", 5000, "Connection dialing timeout, in ms")
	rootCmd.Flags().BoolVarP(&params.debugPtr, "verbose", "v", false, "Print debugging information to the standard error")
	rootCmd.Flags().IntVarP(&params.reportingIntervalPtr, "interval", "i", 1, "Interval, in seconds, between stats updates")
	rootCmd.Flags().BoolVarP(&params.assumeyesPtr, "assume-yes", "y", false, "Force execution without asking for confirmation")
}

func validateRequiredArgs(params *tcpgoonParams, args []string) error {
	if len(args) != 2 {
		return errors.New("Number of required parameters doesn't match")
	}
	params.hostPtr = args[0]
	if port, err := strconv.Atoi(args[1]); err != nil && port <= 0 {
		return errors.New("Port argument is not a valid integer")
	} else {
		params.portPtr = port
	}
	return nil
}

func enableDebugging(params tcpgoonParams) {
	if params.debugPtr {
		debugging.EnableDebug()
	}
}

func autorunValidation(params tcpgoonParams) {
	if !(params.assumeyesPtr || cmdutil.AskForUserConfirmation(params.hostPtr, params.portPtr, params.numberConnectionsPtr)) {
		fmt.Fprintln(debugging.DebugOut, "Execution not approved by the user")
		cmdutil.CloseAbruptly()
	}
}
