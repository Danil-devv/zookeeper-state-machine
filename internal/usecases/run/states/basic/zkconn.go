package basic

import (
	"github.com/go-zookeeper/zk"
	"time"
)

type ZkConn interface {
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	Delete(path string, version int32) error
	Children(path string) ([]string, *zk.Stat, error)
	Exists(path string) (bool, *zk.Stat, error)
	Reconnect(servers []string, timeout time.Duration) error
	Close()
}
