package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dachad/tcpgoon/cmdutil"
	"github.com/dachad/tcpgoon/debugging"
	"github.com/dachad/tcpgoon/mtcpclient"
	"github.com/dachad/tcpgoon/tcpclient"
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

var runCmd = &cobra.Command{
	Use:   "run [flags] <host> <port>",
	Short: "Run tcpgoon test",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := validateRequiredArgs(&params, args); err != nil {
			cmd.Println(cmd.UsageString())
			os.Exit(1)
		}
		enableDebuggingIfFlagSet(params)
		autorunValidation(params)
	},
	Run: func(cmd *cobra.Command, args []string) {
		run(params)
	},
}

func init() {
	runCmd.Flags().IntVarP(&params.numberConnectionsPtr, "connections", "c", 100, "Number of connections you want to open")
	runCmd.Flags().IntVarP(&params.delayPtr, "sleep", "s", 10, "Time you want to sleep between connections, in ms")
	runCmd.Flags().IntVarP(&params.connDialTimeoutPtr, "dial-timeout", "t", 5000, "Connection dialing timeout, in ms")
	runCmd.Flags().BoolVarP(&params.debugPtr, "debug", "d", false, "Print debugging information to the standard error")
	runCmd.Flags().IntVarP(&params.reportingIntervalPtr, "interval", "i", 1, "Interval, in seconds, between stats updates")
	runCmd.Flags().BoolVarP(&params.assumeyesPtr, "assume-yes", "y", false, "Force execution without asking for confirmation")
}

func validateRequiredArgs(params *tcpgoonParams, args []string) error {
	if len(args) != 2 {
		return errors.New("Number of required parameters doesn't match")
	}
	params.hostPtr = args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil && port <= 0 {
		return errors.New("Port argument is not a valid integer")
	}
	params.portPtr = port

	return nil
}

func enableDebuggingIfFlagSet(params tcpgoonParams) {
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

func run(params tcpgoonParams) {
	tcpclient.DefaultDialTimeoutInMs = params.connDialTimeoutPtr

	// TODO: we should decouple the caller from the mtcpclient package (too many structures being moved from
	//  one side to the other.. everything in a single structure, or applying something like the builder pattern,
	//  may help
	connStatusCh, connStatusTracker := mtcpclient.StartBackgroundReporting(params.numberConnectionsPtr, params.reportingIntervalPtr)
	closureCh := mtcpclient.StartBackgroundClosureTrigger(connStatusTracker)
	mtcpclient.MultiTCPConnect(params.numberConnectionsPtr, params.delayPtr, params.hostPtr, params.portPtr, connStatusCh, closureCh)
	fmt.Fprintln(debugging.DebugOut, "Tests execution completed")

	cmdutil.CloseNicely(params.hostPtr, params.portPtr, connStatusTracker)
}
