package cmd

import (
	"log"

	"github.com/artemis19/viz/rpc"
	"github.com/artemis19/viz/server/api"
	"github.com/artemis19/viz/server/config"
	"github.com/artemis19/viz/server/database"
	"github.com/artemis19/viz/server/websocket"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	config.InitializeConfig()
	database.InitializeDB()
	log.Printf("Server is being configured and initialized.\n")
	// Run websocket & rpc in goroutine to avoid blocking functions
	go websocket.Serve()
	go rpc.RunGRPCServer()
	api.Serve()
}
