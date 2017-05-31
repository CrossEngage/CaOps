package gossip

import "github.com/hashicorp/serf/serf"

// UserEventHandler receives the event name, a payload (up to 512 bytes), and must return
// if the event handler processing must break or not, and if necessary, return an error
type UserEventHandler func(event serf.UserEvent) (breakLoop bool, err error)
