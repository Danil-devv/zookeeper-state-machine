package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"time"
)

type Conn struct {
	*zk.Conn
}

func NewConn(servers []string, timeout time.Duration) (*Conn, error) {
	conn, err := connect(servers, timeout)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: conn}, nil
}

func (c *Conn) Reconnect(servers []string, timeout time.Duration) error {
	conn, err := connect(servers, timeout)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func connect(servers []string, timeout time.Duration) (*zk.Conn, error) {
	conn, _, err := zk.Connect(servers, timeout)
	return conn, err
}
