package irc

import (
	"fmt"
)

type Client struct {
	conn     *Conn
	addr     string
	Received chan *Message
	Error    chan error
	quit     chan error
}

func Connect(server string, port int) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", server, port)
	conn, err := Dial(addr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn:     conn,
		addr:     addr,
		Received: make(chan *Message, 32),
		Error:    make(chan error, 1),
	}

	go c.reader()

	go func() {
		err := <-c.quit
		c.quit <- err
		if err != nil {
			c.Error <- err
		}
		close(c.Error)
		close(c.Received)
		c.conn.Close()
	}()

	return c, nil
}

func (c *Client) reader() {
	for {
		m, err := c.conn.ReadMessage()
		if err != nil {
			c.quit <- err
			return
		}

		switch string(m.Command) {
		case "PING":
			if err := c.Send("PONG %s", m.Trailing); err != nil {
				c.quit <- err
				return
			}
		default:
			c.Received <- m
		}
	}
}

func (c *Client) Send(msg string, args ...interface{}) error {
	_, err := fmt.Fprintf(c.conn, msg, args...)
	return err
}

func (c *Client) Close() {
	c.quit <- nil
}
