# Athena

## Arquitecture

* Interacts with Cassandra via Jolokia+JMX
* Uses Cassandra as a backing storage
* Has subcommands to handle bootstrapping separately of serving
* Pools every 30 seconds for new scheduled jobs

* Daemon with HTTP API
* Exports an API `at`-like for scheduling the backup of a keyspace and/or specific tables
* Exports an API `cron`-like for creating recurrent backup jobs
* Does many snapshots at a time
* Uploads the snapshots in many queues
* Controls snapshot allocation
* Watch snapshot disk-usage
* Watch Cassandra cluster state (missing nodes)
* Online DB migration (detect missing keyspace, and creates it)
* Check node status
* Check cluster status


## To-Do

* Better error handling of Jolokia errors
* Tests
* Better code org


## URLs

https://jolokia.org/reference/html/agents.html#agents-jvm