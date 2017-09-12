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

Manages cluster communication, and event handlers

## Scheduler

* Triggers events to Local Agent
* Keeps state on Distributed State Machine