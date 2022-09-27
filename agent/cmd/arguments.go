package cmd

import (
	"log"

	"github.com/artemis19/viz/agent/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:     "agent",
	Version: "0.0.1",
	Short:   "Local traffic capture agent",
	//Run:     run,
}

func Execute() {
	// Do not include "completion" subcommand
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	config.InitializeLogging()

	log.Printf("Agent %v is booting up...", RootCmd.Version)
	defer log.Printf("Agent %v is shutting down...", RootCmd.Version)

	if err := RootCmd.Execute(); err != nil {
		println(err.Error())
	}
}

func initListenCommand() {
	listenCmd := &cobra.Command{
		Run:   Listen,
		Use:   "listen",
		Short: "Listen for traffic",
	}

	// interface, -i
	listenCmd.PersistentFlags().StringP("interface", "i", "", "Network interface to listen on")
	viper.BindPFlag("interface", listenCmd.PersistentFlags().Lookup("interface"))

	// configFile, -c, default is agent.yml
	listenCmd.PersistentFlags().StringP("configFile", "c", "agent.yml", "Location of config file to read from")
	viper.BindPFlag("configFile", listenCmd.PersistentFlags().Lookup("configFile"))

	// outputFile, -o, default is agent.log
	listenCmd.PersistentFlags().StringP("outputFile", "o", "agent.log", "Location of log file output")
	viper.BindPFlag("outputFile", listenCmd.PersistentFlags().Lookup("outputFile"))

	// filters, -f
	listenCmd.PersistentFlags().StringP("filter", "f", "", "Packet capture filters for agent traffic")
	viper.BindPFlag("filter", listenCmd.PersistentFlags().Lookup("filter"))

	// snapLength, -l
	listenCmd.PersistentFlags().StringP("snapLength", "l", "262144", "Snap length setting")
	viper.BindPFlag("snapLength", listenCmd.PersistentFlags().Lookup("snapLength"))

	RootCmd.AddCommand(listenCmd)
}

func initConfigCommand() {
	configCmd := &cobra.Command{
		Run:   Config,
		Use:   "config",
		Short: "Show configuration settings",
	}

	// configFile, -c, default is agent.yml
	configCmd.PersistentFlags().StringP("configFile", "c", "agent.yml", "Location of config file to read from")

	RootCmd.AddCommand(configCmd)
}

func initPingCommand() {
	pingCmd := &cobra.Command{
		Run:   Ping,
		Use:   "ping",
		Short: "Test connectivity to collection server",
	}

	RootCmd.AddCommand(pingCmd)
}

func init() {

	// debug, -d
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Turn debug mode on")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	// serverAddress, -s
	RootCmd.PersistentFlags().StringP("serverAddress", "s", "", "Server to connect to")
	viper.BindPFlag("serverAddress", RootCmd.PersistentFlags().Lookup("serverAddress"))
	// operating system, -p
	RootCmd.PersistentFlags().StringP("operatingSystem", "p", "", "Force operating system for debugging purposes")
	viper.BindPFlag("operatingSystem", RootCmd.PersistentFlags().Lookup("operatingSystem"))

	initListenCommand()
	initConfigCommand()
	initPingCommand()
}
