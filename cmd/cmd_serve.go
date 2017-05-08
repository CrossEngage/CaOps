package cmd

import "github.com/spf13/cobra"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the main daemon",
	Run:   runServeCmd,
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) {

}
