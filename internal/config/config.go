package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ZookeeperServers      []string
	LeaderTimeout         time.Duration
	AttempterTimeout      time.Duration
	FailoverTimeout       time.Duration
	FileDir               string
	StorageCapacity       int
	FailoverAttemptsCount int
}

func GetEnvConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("cannot load .env file: %s\n", err.Error())
	}
	conf := &Config{}
	conf.ZookeeperServers = loadStringSliceEnvVariable(
		"ELECTION_ZK_SERVERS",
		[]string{"zoo1:2181", "zoo2:2182", "zoo3:2183"},
	)
	conf.LeaderTimeout = loadDurationEnvVariable(
		"ELECTION_LEADER_TIMEOUT",
		10*time.Second,
	)
	conf.AttempterTimeout = loadDurationEnvVariable(
		"ELECTION_ATTEMPTER_TIMEOUT",
		2*time.Second,
	)
	conf.FailoverTimeout = loadDurationEnvVariable(
		"ELECTION_FAILOVER_TIMEOUT",
		1*time.Second,
	)
	conf.FileDir = loadStringEnvVariable(
		"ELECTION_FILE_DIR",
		"/tmp/election",
	)
	conf.StorageCapacity = loadIntEnvVariable(
		"ELECTION_STORAGE_CAPACITY",
		10,
	)
	conf.FailoverAttemptsCount = loadIntEnvVariable(
		"ELECTION_FAILOVER_ATTEMPTS_COUNT",
		10,
	)
	return conf
}

func loadStringEnvVariable(name string, defaultValue string) string {
	s, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}

	return s
}

func loadIntEnvVariable(name string, defaultValue int) int {
	s, exists := os.LookupEnv(name)
	res, err := strconv.Atoi(s)
	if !exists || err != nil {
		return defaultValue
	}

	return res
}

func loadStringSliceEnvVariable(name string, defaultValue []string) []string {
	s, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}

	return strings.Split(s, ",")
}

func loadDurationEnvVariable(name string, defaultValue time.Duration) time.Duration {
	s, exists := os.LookupEnv(name)
	d, err := time.ParseDuration(s)
	if !exists || err != nil {
		return defaultValue
	}

	return d
}
