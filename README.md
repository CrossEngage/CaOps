# Athena

## Arquitecture

* Interacts with Cassandra via Jolokia+JMX
* Uses Cassandra as a backing storage
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

* job scheduling at cron and at?
* subcommand to trigger backups (http api)
* Better error handling of Jolokia errors
* Tests
* Better code org
* log to file
* cluster log to separate file (serf)
* glide dep management before 1st release
* signal handling

## Priam

* Token management using SimpleDB
* Support multi-region Cassandra deployment in AWS via public IP.
* Automated security group update in multi-region environment.
* Backup SSTables from local ephemeral disks to S3.
* Uses Snappy compression to compress backup data on the fly.
* Backup throttling
* Pluggable modules for future enhancements (support for multiple data storage).
* APIs to list and restore backup data.
* REST APIs for backup/restore and other operations