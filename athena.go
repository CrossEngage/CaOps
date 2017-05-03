//go:generate make gen_version
package main

import (
	"net/http"
	"os"

	"bitbucket.org/crossengage/athena/cassandra"
	"bitbucket.org/crossengage/athena/cassandra/jolokia"

	"github.com/gocql/gocql"
	logging "github.com/op/go-logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName = "athena"
)

var (
	app = kingpin.New(appName, "A service to backup Cassandra keyspaces to MS Azure")

	debug             = app.Flag("debug", "Enable debugging.").Short('d').Default("false").Bool()
	cqlAddr           = app.Flag("cql", "CQL address of Cassandra node.").Short('c').Default("127.0.0.1:9042").TCP()
	jolokiaURL        = app.Flag("jolokia", "Jolokia URL of Cassandra node.").Short('j').Default("http://127.0.0.1:8778/jolokia").URL()
	protoVer          = app.Flag("protocol-version", "Prefer '3' for C* < 2.2, and '4' for C* >= 2.2, 3.x).").Default("4").Int()
	timeout           = app.Flag("timeout", "Connection timeout.").Short('t').Default("600ms").Duration()
	sockKeepAlive     = app.Flag("keep-alive", "Socket keep-alive interval.").Short('k').Default("0").Duration()
	reconnectInterval = app.Flag("reconnect-interval", "If > 0, attempt to reconnect known DOWN nodes every this").Short('R').Default("0").Duration()

	loggr = logging.MustGetLogger(appName)
)

func init() {
	app.Version(version)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		logging.SetLevel(logging.DEBUG, appName)
		logging.SetFormatter(logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{level} ¶ %{shortfile} ▶ %{message}%{color:reset}`))
	} else {
		logging.SetLevel(logging.INFO, appName)
		logging.SetFormatter(logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{level:-7s} ▶ %{message}%{color:reset}`))
	}
}

func main() {
	cluster := gocql.NewCluster((*cqlAddr).IP.String())
	cluster.Port = (*cqlAddr).Port
	cluster.ProtoVersion = *protoVer
	cluster.Consistency = gocql.ParseConsistency("QUORUM")
	cluster.Timeout = *timeout
	cluster.SocketKeepalive = *sockKeepAlive
	cluster.DisableInitialHostLookup = true
	cluster.ReconnectInterval = *reconnectInterval
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	cluster.Compressor = gocql.SnappyCompressor{}

	session, err := cluster.CreateSession()
	if err != nil {
		loggr.Fatal(err)
	}
	defer session.Close()

	nodetool := cassandra.NodeTool{
		JolokiaClient: jolokia.Client{
			HTTPClient: http.DefaultClient,
			BaseURL:    **jolokiaURL,
		},
	}

	liveNodes, err := nodetool.LiveNodes()
	if err != nil {
		loggr.Fatal(err)
	}

	loggr.Infof("%#v", liveNodes)
	loggr.Info("That's all folks!")
}
