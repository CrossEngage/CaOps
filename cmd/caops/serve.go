package main

import (
	"log"

	"github.com/crossengage/CaOps/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the main daemon",
	Run:   runServeCmd,
}

func init() {
	baseCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) {
	CaOps, err := server.NewCaOps(
		viper.GetString("api.server.bind_addr"),
		viper.GetString("gossip.bind_addr"),
		viper.GetString("gossip.snapshot_path"),
		viper.GetString("cassandra.jolokia_url"),
	)
	if err != nil {
		log.Fatal(err)
	}

	CaOps.Run()
}
