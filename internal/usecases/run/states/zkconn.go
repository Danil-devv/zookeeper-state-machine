package states

import (
	"time"

	"github.com/go-zookeeper/zk"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --output=./mocks --name=ZkConn
type ZkConn interface {
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	Delete(path string, version int32) error
	Children(path string) ([]string, *zk.Stat, error)
	Exists(path string) (bool, *zk.Stat, error)
	Reconnect(servers []string, timeout time.Duration) error
	CheckConnection() bool
	Close()
}
