package main

import (
	"time"
)

type Stats struct {
	/* opaque data */
}

type Server interface {
	ServerID() string
	GetStats() (Stats, error)
}

type Database interface {
	Write(serverID string, stats Stats) error
}

type MonitoringServer struct {
	Servers      []Server
	DB           Database
	pollInterval time.Duration
}

// Monitor runs forever, periodically polling each Server and
// writing its Stats to the Database exactly once per pollInterval.
func (m *MonitoringServer) Monitor() {
	// TODO: implement me!
}
