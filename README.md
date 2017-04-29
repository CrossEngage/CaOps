# Athena

> When Troy fell to the Greeks, Cassandra tried to find a shelter in Athena’s Temple, but she was brutally abducted by Ajax and was brought to Agamemnon as a concubine. Cassandra died in Mycenae, murdered along with Agamemnon by his wife Clytemnestra and her lover Aegisthus.


# Arquitecture

* Connects to Cassandra via Jolokia+JMX

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


# OPT

* Should rename it to cerberus? (something with many heads?)
* Serf internode comm, to know all Athenas are online


# URLs

https://jolokia.org/reference/html/agents.html#agents-jvm