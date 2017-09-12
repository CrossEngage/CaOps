# Architecture

## CaOps

* Integrates all components
* Register cluster event handlers
* Listens for HTTP calls to the API
* Send cluster-wide events
* Displays distributed state

## Cassandra Manager

* Talks to local Cassandra to do operations
* Watches for completion

## Cassandra Jolokia Client

* Talks with Cassandra using Jolokia API

## SnapshotHandler

* Uploads files to remote storage while compressing

## Gossiper

* Sends commands to cluster
* Receives commands from cluster
* Call event handlers
* Accept event handler registration

## Scheduler

* Triggers events to Local Agent
* Keeps state on Distributed State Machine