package cmd

import (
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.org/crossengage/athena/cassandra"
	"bitbucket.org/crossengage/athena/jolokia"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	simpleVersion = false
	versionCmd    = &cobra.Command{
		Use:   "version",
		Short: "Show the version of Athena, Jolokia, and Cassandra",
		Run:   runVersionCmd,
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&simpleVersion, "simple", false, "If set, shows only the version of Athena")
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("%-10s : %s\n", appName, version)
	if simpleVersion {
		return
	}

	jolokiaClient := getJolokiaClient()
	verResp, err := jolokiaClient.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-10s : %s\n", "jolokia", verResp.Value.Agent)

	nodeprobe := cassandra.NewNodeProbe(jolokiaClient)
	cver, err := nodeprobe.StorageService.ReleaseVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-10s : %s\n", "cassandra", cver)
}

func getJolokiaClient() jolokia.Client {
	jolokiaURL, err := url.Parse(viper.GetString("cassandra.jolokia"))
	if err != nil {
		log.Fatal(err)
	}
	return jolokia.Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    *jolokiaURL,
	}
}
