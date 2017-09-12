# Architecture

## CaOps

* Integrates all components

## Local Agent

* Talks to local Cassandra to do operations (snapshotting)
* Watches for completion
* Uploads files to remote storage while compressing

## HTTP Service

* Listens for HTTP calls to the API
* Communicates with Event Dispatcher to send cluster-wide messages
* Displays distributed state (scheduler)

## Broadcaster

* Sends commands to cluster

## Cluster Listener

* Receives commands from cluster
* Talks with Local Agent

## Distributed State Machine

* Uses Raft to keep cluster state (scheduled tasks)

## Scheduler

* Triggers events to Local Agent
* Keeps state on Distributed State Machine

## Cassandra Jolokia Client

* Talks with Cassandra using Jolokia API

## Jolokia Client

* Generic Jolokia Client