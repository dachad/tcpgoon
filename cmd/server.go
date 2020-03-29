package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/dachad/tcpgoon/tcpserver"

	"github.com/spf13/cobra"
)

type TCPServerParams struct {
	port           int
	maxconnections int
	duration       int
}

var tcpserverparams TCPServerParams

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Simple TCP server",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := validateTCPServerArgs(&tcpserverparams); err != nil {
			cmd.Println(err)
			cmd.Println(cmd.UsageString())
			os.Exit(1)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		runTcpgoonServer(tcpserverparams)
		return nil
	},
}

func init() {
	serverCmd.Flags().IntVarP(&tcpserverparams.port, "port", "p", 54321, "TCP listening port, from 1024 to 65535")
	serverCmd.Flags().IntVarP(&tcpserverparams.maxconnections, "maxconnections", "m", 10, "How many total connections we will accept")
	serverCmd.Flags().IntVarP(&tcpserverparams.duration, "duration", "d", 30, "Running time before dropping")
}

func validateTCPServerArgs(params *TCPServerParams) error {
	if params.port < 1024 || params.port > 65535 {
		return errors.New(strconv.Itoa(params.port) + " is not a valid TCP port number for the server")
	}

	if params.maxconnections < 0 {
		return errors.New("Max Connections argument should be a positive integer")
	}

	if params.duration <= 0 {
		return errors.New("Duration argument should be a positive integer")
	}

	return nil
}

func runTcpgoonServer(params TCPServerParams) {
	fmt.Println("Running the simple TCP server in port", params.port, "up to", params.maxconnections, "connections or", params.duration, "seconds, what happens first")

	dispatcher := &tcpserver.Dispatcher{
		Handlers: make(map[string]*tcpserver.Handler),
		Lock:     sync.RWMutex{},
	}

	var end_waiter sync.WaitGroup
	end_waiter.Add(1)

	runTCPServer := func() {
		fmt.Println("Starting TCP server")
		if err := dispatcher.ListenHandlersComplete(params.port, params.maxconnections, params.duration, &end_waiter); err != nil {
			fmt.Println("Could not start the TCP server", err)
			return
		}
	}
	go runTCPServer()

	WaitForCtrlC(&end_waiter)

	end_waiter.Wait()
}

func WaitForCtrlC(end_waiter *sync.WaitGroup) {
	signal_channel := make(chan os.Signal, 1)
	fmt.Printf("Press Ctrl+C to end\n")
	signal.Notify(signal_channel, os.Interrupt)

	go func() {
		<-signal_channel
		fmt.Println()
		end_waiter.Done()
	}()
}
