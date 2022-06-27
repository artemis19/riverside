package cmd

import (
	"github.com/artemis19/viz/server/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var RootCmd = &cobra.Command{
	Use:     "server",
	Version: "0.0.1",
	Short:   "Server to store agent traffic",
	Run:     Run,
}

func Execute() {
	// Do not include "completion" subcommand
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	config.InitializeLogging()

	log.Printf("Server %v is booting up...", RootCmd.Version)
	defer log.Printf("Server %v is shutting down...", RootCmd.Version)

	if err := RootCmd.Execute(); err != nil {
		println(err.Error())
	}
}

func initConfigCommand() {
	configCmd := &cobra.Command{
		Run:   Config,
		Use:   "config",
		Short: "Show configuration settings",
	}

	// configFile, -c, default is agent.yml
	configCmd.PersistentFlags().StringP("configFile", "c", "server.yml", "Location of config file to read from")

	RootCmd.AddCommand(configCmd)
}

func init() {
	// configFile, -c, default is agent.yml
	RootCmd.PersistentFlags().StringP("configFile", "c", "server.yml", "Location of config file to read from")
	viper.BindPFlag("configFile", RootCmd.PersistentFlags().Lookup("configFile"))

	// outputFile, -o, default is agent.log
	RootCmd.PersistentFlags().StringP("outputFile", "o", "server.log", "Location of log file output")
	viper.BindPFlag("outputFile", RootCmd.PersistentFlags().Lookup("outputFile"))

	// port, -l
	RootCmd.PersistentFlags().StringP("port", "l", "1533", "Port for server to listen on")
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))

	// websocket port, -w
	RootCmd.PersistentFlags().StringP("websocketPort", "w", "8000", "Websocker server listening port")
	viper.BindPFlag("websocketPort", RootCmd.PersistentFlags().Lookup("websocketPort"))

	// DB IP address/host, -a
	RootCmd.PersistentFlags().StringP("dbAddress", "a", "localhost", "Postgres database address")
	viper.BindPFlag("dbAddress", RootCmd.PersistentFlags().Lookup("dbAddress"))

	// DB port, -b
	RootCmd.PersistentFlags().StringP("dbPort", "b", "5432", "Postgres database port")
	viper.BindPFlag("dbPort", RootCmd.PersistentFlags().Lookup("dbPort"))

	// DB username, -u
	RootCmd.PersistentFlags().StringP("dbUsername", "u", "viz", "Postgres database user")
	viper.BindPFlag("dbUsername", RootCmd.PersistentFlags().Lookup("dbUsername"))

	// DB password, -p
	RootCmd.PersistentFlags().StringP("dbPassword", "p", "mysecretpassword", "Postgres database password")
	viper.BindPFlag("dbPassword", RootCmd.PersistentFlags().Lookup("dbPassword"))

	// DB SSL, -s
	RootCmd.PersistentFlags().BoolP("dbSSL", "s", false, "Turn Postgres SSL mode on")
	viper.BindPFlag("dbSSL", RootCmd.PersistentFlags().Lookup("dbSSL"))

	// debug, -d
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Turn debug mode on")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	initConfigCommand()
}
