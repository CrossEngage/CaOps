package main

// TODO - move this to integration tests

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/gocql/gocql"
	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pb "gopkg.in/cheggaaa/pb.v1"
	inf "gopkg.in/inf.v0"
)

var (
	fixtureCmd = &cobra.Command{
		Use:   "fixture",
		Short: "Create keyspaces with fake data, for testing",
		Run:   runFixtureCmd,
	}
	cassandraAddr string
	cassandraUser string
	cassandraPass string
	keyspaceName  string
	numProducts   int
	numUsers      int
)

func init() {
	baseCmd.AddCommand(fixtureCmd)
	fixtureCmd.Flags().StringVar(&cassandraAddr, "host", "127.0.0.1:9042", "Cassandra CQL host:port")
	fixtureCmd.Flags().StringVar(&cassandraUser, "user", "cassandra", "Cassandra Username")
	fixtureCmd.Flags().StringVar(&cassandraPass, "pass", "cassandra", "Cassandra Password")
	fixtureCmd.Flags().StringVar(&keyspaceName, "keyspace", "company_xyz", "The keyspace name to create and fill")
	fixtureCmd.Flags().IntVar(&numProducts, "num-products", 100, "Number of products to create")
	fixtureCmd.Flags().IntVar(&numUsers, "num-users", 10000, "Number of users to create")
}

func logFatal(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func runFixtureCmd(cmd *cobra.Command, args []string) {
	cluster := getCassandraClusterConfig()
	session, err := cluster.CreateSession()
	logrus.Infof("Connecting to %s", cassandraAddr)
	logFatal(err)
	logFatal(createKeyspace(session, keyspaceName))
	logFatal(createCountersTable(session, keyspaceName))
	logFatal(createProductsTable(session, keyspaceName))
	logFatal(fillProductsTable(session, keyspaceName, numProducts))
	logFatal(createUsersTable(session, keyspaceName))
	logFatal(fillUsersTable(session, keyspaceName, numUsers))
	defer session.Close()
}

func getCassandraClusterConfig() *gocql.ClusterConfig {
	cqlAddr, err := net.ResolveTCPAddr("tcp", cassandraAddr)
	logFatal(err)
	config := gocql.NewCluster(cqlAddr.IP.String())
	config.Port = cqlAddr.Port
	config.Authenticator = &gocql.PasswordAuthenticator{Username: cassandraUser, Password: cassandraPass}
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
	logrus.Infof("Creating keyspace %s at %s", keyspaceName, cassandraAddr)
	query := `
		CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION =
		{'class':'SimpleStrategy','replication_factor':3}`
	if err := session.Query(fmt.Sprintf(query, keyspace)).Exec(); err != nil {
		return err
	}
	return nil
}

func createCountersTable(session *gocql.Session, keyspace string) error {
	logrus.Infof("Creating counters table at %s@%s", keyspaceName, cassandraAddr)
	query := `
		CREATE TABLE IF NOT EXISTS %s.counters (
			id    varchar PRIMARY KEY,
			count counter
		) WITH comment='Counters'`
	if err := session.Query(fmt.Sprintf(query, keyspace)).Exec(); err != nil {
		return err
	}
	return nil
}

func createProductsTable(session *gocql.Session, keyspace string) error {
	logrus.Infof("Creating products table at %s@%s", keyspaceName, cassandraAddr)
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
	logrus.Infof("Filling products table at %s@%s", keyspaceName, cassandraAddr)
	bar := pb.StartNew(qtd)
	for i := 0; i < qtd; i++ {
		bar.Increment()
		query := `
			INSERT INTO %s.products (sku, brand, name, model, price)
			VALUES (?, ?, ?, ?, ?) IF NOT EXISTS`
		if err := session.Query(fmt.Sprintf(query, keyspace),
			uuid.NewV4().String(), fake.Brand(), fake.ProductName(), fake.Model(),
			inf.NewDec(rand.Int63()+1, inf.Scale(rand.Int31()))).Exec(); err != nil {
			return err
		}
		if err := incrementCounter(session, keyspace, "Products"); err != nil {
			return err
		}
	}
	bar.FinishPrint("Done.")
	return nil
}

func createUsersTable(session *gocql.Session, keyspace string) error {
	logrus.Infof("Creating users table at %s@%s", keyspaceName, cassandraAddr)
	query := `
		CREATE TABLE IF NOT EXISTS %s.users (
			id        uuid PRIMARY KEY,
			username  varchar,
			email     varchar,
			full_name varchar,
			gender    varchar,
		) WITH comment='Users'`
	if err := session.Query(fmt.Sprintf(query, keyspace)).Exec(); err != nil {
		return err
	}
	return nil
}

func fillUsersTable(session *gocql.Session, keyspace string, qtd int) error {
	logrus.Infof("Filling users table at %s@%s", keyspaceName, cassandraAddr)
	bar := pb.StartNew(qtd)
	for i := 0; i < qtd; i++ {
		bar.Increment()
		query := `
			INSERT INTO %s.users (id, username, email, full_name, gender)
			VALUES (?, ?, ?, ?, ?) IF NOT EXISTS`
		if err := session.Query(fmt.Sprintf(query, keyspace),
			uuid.NewV4().String(), fake.UserName(), fake.EmailAddress(),
			fake.FullName(), fake.Gender()).Exec(); err != nil {
			return err
		}
		if err := incrementCounter(session, keyspace, "Users"); err != nil {
			return err
		}
	}
	bar.FinishPrint("Done.")
	return nil
}

func incrementCounter(session *gocql.Session, keyspace, counter string) error {
	query := `UPDATE %s.counters SET count = count + 1 WHERE id = ?`
	if err := session.Query(fmt.Sprintf(query, keyspace), counter).Exec(); err != nil {
		return err
	}
	return nil
}
