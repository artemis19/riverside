package cmd

import (
	"github.com/artemis19/viz/rpc"
	"github.com/artemis19/viz/server/config"
	"github.com/artemis19/viz/server/database"
	"github.com/artemis19/viz/server/websocket"
	"github.com/spf13/cobra"
	"log"
)

func Run(cmd *cobra.Command, args []string) {
	config.InitializeConfig()
	database.InitializeDB()
	log.Printf("Server is being configured and initialized.\n")
	// Run websocket in goroutine
	go websocket.Serve()
	rpc.RunGRPCServer()
}
