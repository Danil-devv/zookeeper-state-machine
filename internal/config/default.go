package config

import "time"

var DefaultConfigValue = Config{
	ZookeeperServers:      []string{"zoo1:2181", "zoo2:2182", "zoo3:2183"},
	LeaderTimeout:         10 * time.Second,
	AttempterTimeout:      2 * time.Second,
	FailoverTimeout:       1 * time.Second,
	FileDir:               "/tmp/election",
	StorageCapacity:       10,
	FailoverAttemptsCount: 10,
}
