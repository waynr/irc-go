package irc

import (
	"net"
)

type Connection struct {
	conn     net.Conn
}

func Dial(address string) (*Connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	c := &Connection{conn: conn}
	return c, nil
}

func (c *Connection) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}

func (c *Connection) Read(p []byte) (int, error) {
	return c.conn.Read(p)
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
