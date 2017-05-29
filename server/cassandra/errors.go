package cassandra

import "errors"

var (
	ErrUnreachableCassandraNodes = errors.New("Unreachable Cassandra nodes")
	ErrJoiningCassandraNodes     = errors.New("Joining Cassandra nodes")
	ErrLeavingCassandraNodes     = errors.New("Leaving Cassandra nodes")
	ErrMovingCassandraNodes      = errors.New("Moving Cassandra nodes")
)
