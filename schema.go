package main

var (
	migrations = map[int]string{

		20170502000: `CREATE KEYSPACE IF NOT EXISTS athena`,

		20170502001: `CREATE TABLE IF NOT EXISTS athena.schema_version (
						version bigint,
						PRIMARY KEY (version)
					  ) WITH CLUSTERING ORDER BY (version DESC) 
					  AND comment='DB Schema version'`,

		20170502010: `CREATE TABLE IF NOT EXISTS athena.recurring (
						id         uuid,
						every      int,
						period     text,
						at_hour    int,
						at_minute  int,
						keyspace   text,
						table      text,
						PRIMARY KEY (id),
					  ) WITH comment=''`,

		20170502020: `CREATE TABLE IF NOT EXISTS athena.one_time (
						date timestamp, 
						when timeuuid,
						keyspace text,
						table text,
						PRIMARY KEY (date),
			          ) WITH CLUSTERING ORDER BY (when DESC) 
					  AND comment='Schedule one time backups' 
					  AND default_time_to_live=605000`,
	}
)
