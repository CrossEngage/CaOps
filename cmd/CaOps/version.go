package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	complete   = false
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version of CaOps, Jolokia, and Cassandra",
		Run:   runVersionCmd,
	}
)

func init() {
	baseCmd.AddCommand(versionCmd)
	// TODO add this flag when the following to-do is done
	// versionCmd.Flags().BoolVar(&complete, "complete", false,
	// "If set, shows the version of the Jolokia Agent and Cassandra this agent is connected to")
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	log.Printf("%s: %s\n", appName, version)
	if complete {
		// TODO initialize Cassandra Manager and show Jolokia Agent and Cassandra version info
		return
	}
}
