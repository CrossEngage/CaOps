package main

import (
	"log"

	"bitbucket.org/crossengage/athena/server"
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
	athena, err := server.NewAthena(
		viper.GetString("api.server.bind_addr"),
		viper.GetString("gossip.bind_addr"),
		viper.GetString("gossip.snapshot_path"),
		viper.GetString("cassandra.jolokia_url"),
	)
	if err != nil {
		log.Fatal(err)
	}

	athena.Run()
}
