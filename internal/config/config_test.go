package config

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testEnvVar struct {
	value     string
	isCorrect bool
}

func setEnvConfig(conf map[string]testEnvVar, t *testing.T) {
	for k, v := range conf {
		err := os.Setenv(k, v.value)
		assert.NoError(t, err)
	}
}

func validateConfig(expected map[string]testEnvVar, actual *Config, t *testing.T) {
	v, ok := expected["ELECTION_ZK_SERVERS"]
	if !ok || !v.isCorrect {
		assert.Equal(t, DefaultConfigValue.ZookeeperServers, actual.ZookeeperServers)
	} else {
		assert.Equal(t, strings.Split(v.value, ","), actual.ZookeeperServers)
	}

	v, ok = expected["ELECTION_LEADER_TIMEOUT"]
	if !ok || !v.isCorrect {
		assert.Equal(t, DefaultConfigValue.LeaderTimeout, actual.LeaderTimeout)
	} else {
		s, err := time.ParseDuration(v.value)
		assert.NoError(t, err)
		assert.Equal(t, s, actual.LeaderTimeout)
	}

	v, ok = expected["ELECTION_ATTEMPTER_TIMEOUT"]
	if !ok || !v.isCorrect {
		assert.Equal(t, DefaultConfigValue.AttempterTimeout, actual.AttempterTimeout)
	} else {
		s, err := time.ParseDuration(v.value)
		assert.NoError(t, err)
		assert.Equal(t, s, actual.AttempterTimeout)
	}

	v, ok = expected["ELECTION_FAILOVER_TIMEOUT"]
	if !ok || !v.isCorrect {
		assert.Equal(t, DefaultConfigValue.FailoverTimeout, actual.FailoverTimeout)
	} else {
		s, err := time.ParseDuration(v.value)
		assert.NoError(t, err)
		assert.Equal(t, s, actual.FailoverTimeout)
	}

	v, ok = expected["ELECTION_FILE_DIR"]
	if !ok || !v.isCorrect {
		assert.Equal(t, DefaultConfigValue.FileDir, actual.FileDir)
	} else {
		assert.Equal(t, v.value, actual.FileDir)
	}

	v, ok = expected["ELECTION_STORAGE_CAPACITY"]
	if !ok || !v.isCorrect {
		assert.Equal(t, DefaultConfigValue.StorageCapacity, actual.StorageCapacity)
	} else {
		x, err := strconv.Atoi(v.value)
		assert.NoError(t, err)
		assert.Equal(t, x, actual.StorageCapacity)
	}

	v, ok = expected["ELECTION_FAILOVER_ATTEMPTS_COUNT"]
	if !ok || !v.isCorrect {
		assert.Equal(t, DefaultConfigValue.FailoverAttemptsCount, actual.FailoverAttemptsCount)
	} else {
		x, err := strconv.Atoi(v.value)
		assert.NoError(t, err)
		assert.Equal(t, x, actual.FailoverAttemptsCount)
	}
}

func TestGetEnvConfig_EmptyEnvFile(t *testing.T) {
	conf := GetEnvConfig()
	assert.Equal(t, DefaultConfigValue, *conf)
}

func TestGetEnvConfig_CorrectFullEnvFile(t *testing.T) {
	tests := []map[string]testEnvVar{
		{
			"ELECTION_ZK_SERVERS":              {"zk1:2234,zk2:2345", true},
			"ELECTION_LEADER_TIMEOUT":          {"10s", true},
			"ELECTION_ATTEMPTER_TIMEOUT":       {"3s", true},
			"ELECTION_FAILOVER_TIMEOUT":        {"5s", true},
			"ELECTION_FILE_DIR":                {"/path/to/dir", true},
			"ELECTION_STORAGE_CAPACITY":        {"10", true},
			"ELECTION_FAILOVER_ATTEMPTS_COUNT": {"10", true},
		},
		{
			"ELECTION_ZK_SERVERS":              {"zookeeper_server1:2120", true},
			"ELECTION_LEADER_TIMEOUT":          {"5s", true},
			"ELECTION_ATTEMPTER_TIMEOUT":       {"3s", true},
			"ELECTION_FAILOVER_TIMEOUT":        {"10s", true},
			"ELECTION_FILE_DIR":                {"/test", true},
			"ELECTION_STORAGE_CAPACITY":        {"5", true},
			"ELECTION_FAILOVER_ATTEMPTS_COUNT": {"2", true},
		},
	}

	for _, tc := range tests {
		setEnvConfig(tc, t)
		conf := GetEnvConfig()
		validateConfig(tc, conf, t)
		os.Clearenv()
	}
}

func TestGetEnvConfig_CorrectNotFullEnvFile(t *testing.T) {
	tests := []map[string]testEnvVar{
		{
			"ELECTION_ZK_SERVERS":       {"zk1:2234,zk2:2345", true},
			"ELECTION_LEADER_TIMEOUT":   {"10s", true},
			"ELECTION_FAILOVER_TIMEOUT": {"5s", true},
			"ELECTION_STORAGE_CAPACITY": {"10", true},
		},
		{
			"ELECTION_ATTEMPTER_TIMEOUT":       {"3s", true},
			"ELECTION_FAILOVER_TIMEOUT":        {"10s", true},
			"ELECTION_FILE_DIR":                {"/test", true},
			"ELECTION_FAILOVER_ATTEMPTS_COUNT": {"2", true},
		},
	}

	for _, tc := range tests {
		setEnvConfig(tc, t)
		conf := GetEnvConfig()
		validateConfig(tc, conf, t)
		os.Clearenv()
	}
}

func TestGetEnvConfig_IncorrectEnvFile(t *testing.T) {
	tests := []map[string]testEnvVar{
		{
			"ELECTION_ZK_SERVERS":              {"zookeeper_server1:2120", true},
			"ELECTION_ATTEMPTER_TIMEOUT":       {"notTimeDuration", false},
			"ELECTION_FAILOVER_TIMEOUT":        {"5s", true},
			"ELECTION_FILE_DIR":                {"th", true},
			"ELECTION_STORAGE_CAPACITY":        {"notNumber", false},
			"ELECTION_FAILOVER_ATTEMPTS_COUNT": {"10", true},
		},
	}

	for _, tc := range tests {
		setEnvConfig(tc, t)
		conf := GetEnvConfig()
		validateConfig(tc, conf, t)
		os.Clearenv()
	}
}
