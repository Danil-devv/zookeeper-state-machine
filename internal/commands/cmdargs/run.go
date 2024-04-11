package cmdargs

import "time"

type RunArgs struct {
	ZookeeperServers      []string
	LeaderTimeout         time.Duration
	AttempterTimeout      time.Duration
	FailoverTimeout       time.Duration
	FileDir               string
	StorageCapacity       int
	FailoverAttemptsCount int
}
