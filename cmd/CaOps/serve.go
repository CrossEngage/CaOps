package main

import (
	"github.com/CrossEngage/CaOps/internal/server"
	"github.com/sirupsen/logrus"
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
		logrus.Fatal(err)
	}

	CaOps.Run()
}
