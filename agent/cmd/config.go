package cmd

import (
	"fmt"
	"github.com/artemis19/viz/agent/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func Config(cmd *cobra.Command, args []string) {
	// Set defaults & start logging
	config.InitializeConfig()

	log.Printf("Tasked to view config info...\n")
	defer log.Printf("Done viewing config info.\n")

	if len(args) == 0 {
		for k, v := range viper.AllSettings() {
			fmt.Printf("%v = %v\n", k, v)
		}
	} else {
		for value := range args {
			fmt.Printf("%v = %v\n", value, viper.Get("value"))
		}
	}
}
