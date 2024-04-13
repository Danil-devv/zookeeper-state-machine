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
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("cannot load .env file: %s\n", err.Error())
	}
	conf := &Config{}
	conf.ZookeeperServers = loadStringSliceEnvVariable(
		"ELECTION_ZK_SERVERS",
		DefaultConfigValue.ZookeeperServers,
	)
	conf.LeaderTimeout = loadDurationEnvVariable(
		"ELECTION_LEADER_TIMEOUT",
		DefaultConfigValue.LeaderTimeout,
	)
	conf.AttempterTimeout = loadDurationEnvVariable(
		"ELECTION_ATTEMPTER_TIMEOUT",
		DefaultConfigValue.AttempterTimeout,
	)
	conf.FailoverTimeout = loadDurationEnvVariable(
		"ELECTION_FAILOVER_TIMEOUT",
		DefaultConfigValue.FailoverTimeout,
	)
	conf.FileDir = loadStringEnvVariable(
		"ELECTION_FILE_DIR",
		DefaultConfigValue.FileDir,
	)
	conf.StorageCapacity = loadIntEnvVariable(
		"ELECTION_STORAGE_CAPACITY",
		DefaultConfigValue.StorageCapacity,
	)
	conf.FailoverAttemptsCount = loadIntEnvVariable(
		"ELECTION_FAILOVER_ATTEMPTS_COUNT",
		DefaultConfigValue.FailoverAttemptsCount,
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
