package schema

import (
	"time"
)

const (
	keyspaceName = `athena`
	defaultTTL   = 30 * 24 * time.Hour
)

var (
	bootstrap = []string{
		`CREATE KEYSPACE IF NOT EXISTS athena`,
		`CREATE TABLE IF NOT EXISTS athena.schema_version (
			version bigint,
			applied timeuuid,
			PRIMARY KEY (version)
		) WITH CLUSTERING ORDER BY (version DESC) 
		AND comment='DB Schema version'`,
	}

	migrations = map[int]string{
		20170502010: `CREATE TABLE IF NOT EXISTS athena.recurring_job (
						id         uuid,
						created	   timeuuid,
						interval   bigint,
						unit       varchar, 
						at_time    varchar,
						last_run   timeuuid,
						next_run   timeuuid,
						tables     text,
						PRIMARY KEY (id),
					  ) WITH comment=''`,

		20170502020: `CREATE TABLE IF NOT EXISTS athena.one_time_job (
						id         uuid,
						created	   timeuuid,
						date       timestamp,
						when	   timeuuid,
						keyspace   text,
						table      text,
						PRIMARY KEY (id),
			          ) WITH CLUSTERING ORDER BY (when DESC) 
					  AND comment='Schedule one time backups' 
					  AND default_time_to_live=2592000`,
	}
)
