package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"

	pb "gopkg.in/cheggaaa/pb.v1"
	inf "gopkg.in/inf.v0"

	"time"

	"github.com/gocql/gocql"
	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var (
	fixtureCmd = &cobra.Command{
		Use:   "fixture",
		Short: "Create keyspaces with fake data, for testing",
		Run:   runFixtureCmd,
	}
	cassandraAddr string
	keyspaceName  string
	numProducts   int
)

func init() {
	baseCmd.AddCommand(fixtureCmd)
	fixtureCmd.Flags().StringVar(&cassandraAddr, "host", "127.0.0.1:9042", "Cassandra CQL host:port")
	fixtureCmd.Flags().StringVar(&keyspaceName, "keyspace", "company_xyz", "The keyspace name to create and fill")
	fixtureCmd.Flags().IntVar(&numProducts, "num-products", 10000, "Number of products to create")
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func runFixtureCmd(cmd *cobra.Command, args []string) {
	cluster := getCassandraClusterConfig()
	session, err := cluster.CreateSession()
	log.Printf("Connecting to %s", cassandraAddr)
	logFatal(err)
	logFatal(createKeyspace(session, keyspaceName))
	logFatal(createProductsTable(session, keyspaceName))
	logFatal(fillProductsTable(session, keyspaceName, numProducts))
	defer session.Close()
}

func getCassandraClusterConfig() *gocql.ClusterConfig {
	cqlAddr, err := net.ResolveTCPAddr("tcp", cassandraAddr)
	if err != nil {
		log.Fatal(err)
	}
	config := gocql.NewCluster(cqlAddr.IP.String())
	config.Port = cqlAddr.Port
	config.ProtoVersion = 4
	config.Consistency = gocql.Quorum
	config.Timeout = 5 * time.Minute
	config.SocketKeepalive = 10 * time.Minute
	config.DisableInitialHostLookup = true
	config.ReconnectInterval = 11 * time.Minute
	config.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	config.MaxWaitSchemaAgreement = 5 * time.Minute
	config.Compressor = gocql.SnappyCompressor{}
	return config
}

func createKeyspace(session *gocql.Session, keyspace string) error {
	log.Printf("Creating keyspace %s at %s", keyspaceName, cassandraAddr)
	query := `
		CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = 
		{'class':'SimpleStrategy','replication_factor':3}`
	if err := session.Query(fmt.Sprintf(query, keyspace)).Exec(); err != nil {
		return err
	}
	return nil
}

func createProductsTable(session *gocql.Session, keyspace string) error {
	log.Printf("Creating products table at %s@%s", keyspaceName, cassandraAddr)
	query := `
		CREATE TABLE IF NOT EXISTS %s.products (
			sku   uuid PRIMARY KEY, 
			brand varchar,
			name  varchar,
			model varchar,
			price decimal,
		) WITH comment='Products'`
	if err := session.Query(fmt.Sprintf(query, keyspace)).Exec(); err != nil {
		return err
	}
	return nil
}

func fillProductsTable(session *gocql.Session, keyspace string, qtd int) error {
	log.Printf("Filling products table at %s@%s", keyspaceName, cassandraAddr)
	bar := pb.StartNew(qtd)
	for i := 0; i < qtd; i++ {
		bar.Increment()
		query := `
			INSERT INTO %s.products (sku, brand, name, model, price)
			VALUES (?, ?, ?, ?, ?) IF NOT EXISTS`
		if err := session.Query(fmt.Sprintf(query, keyspace),
			uuid.NewV4().String(), fake.Brand(), fake.ProductName(), fake.Model(),
			inf.NewDec(rand.Int63()+1, inf.Scale(rand.Int31())),
		).Exec(); err != nil {
			return err
		}
	}
	bar.FinishPrint("Done.")
	return nil
}
