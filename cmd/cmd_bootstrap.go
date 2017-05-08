package cmd

import (
	"net"

	"github.com/gocql/gocql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstraps Athena's keyspace on Cassandra",
	Long: `This sub-command creates Athena's keyspace and tables that are used to store 
	schedule information for the daemon.`,
	Run: runBootstrapCmd,
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
}

func runBootstrapCmd(cmd *cobra.Command, args []string) {
	cluster := getCassandraClusterConfig()
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	if err := bootstrapDatabase(session); err != nil {
		log.Fatal(err)
	}
	log.Info("Bootstrap finished.")
}

func getCassandraClusterConfig() *gocql.ClusterConfig {
	cqlAddr, err := net.ResolveTCPAddr("tcp", viper.GetString("cassandra.cql"))
	if err != nil {
		log.Fatal(err)
	}
	config := gocql.NewCluster(cqlAddr.IP.String())
	config.Port = cqlAddr.Port
	config.ProtoVersion = viper.GetInt("cassandra.protocol")
	config.Consistency = gocql.ParseConsistency("QUORUM")
	config.Timeout = viper.GetDuration("cassandra.timeout")
	config.SocketKeepalive = viper.GetDuration("cassandra.keep-alive")
	config.DisableInitialHostLookup = true
	config.ReconnectInterval = viper.GetDuration("cassandra.reconnect-interval")
	config.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: viper.GetInt("cassandra.retries")}
	config.MaxWaitSchemaAgreement = viper.GetDuration("cassandra.max-wait-schema-agreement")
	config.Compressor = gocql.SnappyCompressor{}
	return config
}
