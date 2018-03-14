package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/dachad/tcpgoon/cmdutil"
	"github.com/dachad/tcpgoon/debugging"
	"github.com/dachad/tcpgoon/mtcpclient"
	"github.com/dachad/tcpgoon/tcpclient"
	"github.com/spf13/cobra"
)

type tcpgoonParams struct {
	target            string
	targetip          string
	port              int
	numberConnections int
	delay             int
	connDialTimeout   int
	debug             bool
	reportingInterval int
	assumeyes         bool
}

var params tcpgoonParams

var runCmd = &cobra.Command{
	Use:   "run [flags] <host> <port>",
	Short: "Run tcpgoon test",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := validateRequiredArgs(&params, args); err != nil {
			cmd.Println(err)
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
	runCmd.Flags().IntVarP(&params.numberConnections, "connections", "c", 100, "Number of connections you want to open")
	runCmd.Flags().IntVarP(&params.delay, "sleep", "s", 10, "Time you want to sleep between connections, in ms")
	runCmd.Flags().IntVarP(&params.connDialTimeout, "dial-timeout", "t", 5000, "Connection dialing timeout, in ms")
	runCmd.Flags().BoolVarP(&params.debug, "debug", "d", false, "Print debugging information to the standard error")
	runCmd.Flags().IntVarP(&params.reportingInterval, "interval", "i", 1, "Interval, in seconds, between stats updates")
	runCmd.Flags().BoolVarP(&params.assumeyes, "assume-yes", "y", false, "Force execution without asking for confirmation")
}

func validateRequiredArgs(params *tcpgoonParams, args []string) error {
	if len(args) != 2 {
		return errors.New("Number of required parameters doesn't match")
	}
	params.target = args[0]
	addrs, err := net.LookupIP(params.target)
	if err != nil || len(addrs) == 0 {
		return errors.New("Domain name not resolvable")
	}

	params.targetip = addrs[0].String()
	fmt.Fprintln(debugging.DebugOut, "TCPGOON target: Hostname(", params.target, "), IP (", params.targetip, ")")

	port, err := strconv.Atoi(args[1])
	if err != nil && port <= 0 {
		return errors.New("Port argument is not a valid integer")
	}
	params.port = port

	return nil
}

func enableDebuggingIfFlagSet(params tcpgoonParams) {
	if params.debug {
		debugging.EnableDebug()
	}
}

func autorunValidation(params tcpgoonParams) {
	if !(params.assumeyes || cmdutil.AskForUserConfirmation(params.target, params.port, params.numberConnections)) {
		fmt.Fprintln(debugging.DebugOut, "Execution not approved by the user")
		cmdutil.CloseAbruptly()
	}
}

func run(params tcpgoonParams) {
	tcpclient.DefaultDialTimeoutInMs = params.connDialTimeout

	// TODO: we should decouple the caller from the mtcpclient package (too many structures being moved from
	//  one side to the other.. everything in a single structure, or applying something like the builder pattern,
	//  may help
	connStatusCh, connStatusTracker := mtcpclient.StartBackgroundReporting(params.numberConnections, params.reportingInterval)
	closureCh := mtcpclient.StartBackgroundClosureTrigger(*connStatusTracker)
	mtcpclient.MultiTCPConnect(params.numberConnections, params.delay, params.target, params.port, connStatusCh, closureCh)
	fmt.Fprintln(debugging.DebugOut, "Tests execution completed")

	cmdutil.CloseNicely(params.targetip, params.target, params.port, *connStatusTracker)
}
