package cache

import "time"

type Topology int

const (
	Standalone Topology = iota
	Cluster
	Sentinel
)

// Config set config for cache
type Config struct {
	MasterName   string        `json:"master_name"`
	Servers      []string      `json:"servers"`
	Timeout      time.Duration `json:"timeout"`
	AuthPass     string        `json:"auth_pass"`
	Topology     Topology      `json:"topology"`
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
}
