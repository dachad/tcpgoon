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
	"github.com/dachad/tcpgoon/docker"

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
	docker			  bool
	containerID		  string
}

var params tcpgoonParams

var runCmd = &cobra.Command{
	Use:   "run [flags] <host> <port>",
	Short: "Run tcpgoon test",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		enableDebuggingIfFlagSet(params)
		if err := validateRequiredArgs(&params, args); err != nil {
			cmd.Println(err)
			cmd.Println(cmd.UsageString())
			os.Exit(1)
		}
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
	runCmd.Flags().BoolVarP(&params.docker, "docker", "", false, "Running against a docker image. " +
	"Use the image reference instead of the <host>, and a custom mapped port if the first exposed one is not ok (optional")
}

func validateRequiredArgs(params *tcpgoonParams, args []string) error {
	var minNumberRequiredArgs int
	var maxNumberRequiredArgs int

	if params.docker {
		minNumberRequiredArgs = 1
		maxNumberRequiredArgs = 2
	} else {
		minNumberRequiredArgs = 2
		maxNumberRequiredArgs = 2
	}

	if len(args) < minNumberRequiredArgs || len(args) > maxNumberRequiredArgs{
		return errors.New("Number of required parameters doesn't match")
	}

	if params.docker {
		dockerImage := args[0]

		fmt.Println("Starting Docker test for", dockerImage)

		var port int
		if len(args) == 2 {

			portInt, err := strconv.Atoi(args[1])
			port = portInt

			if err != nil && port <= 0 {
				return errors.New("Port argument is not a valid integer")
			}
		} else {
			port = 0
		}

		target, containerID := docker.DownloadAndRun(dockerImage, port)
		params.containerID = containerID
		params.target = target.IP
		params.port = target.Port

	} else {
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
	}

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
		if params.docker{
			docker.Stop(params.containerID)
		}
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
	if params.docker{
		go func() {
			<-closureCh
			docker.Stop(params.containerID)
		}()
	}
	mtcpclient.MultiTCPConnect(params.numberConnections, params.delay, params.target, params.port, connStatusCh, closureCh)
	fmt.Fprintln(debugging.DebugOut, "Tests execution completed")

	cmdutil.CloseNicely(params.targetip, params.target, params.port, *connStatusTracker)
}
