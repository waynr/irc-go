package irc

import (
	"bufio"
	"bytes"
	"net"
)

type Conn struct {
	conn net.Conn
}

func Dial(addr string) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Conn{conn}, nil
}

func (c *Conn) ReadLine() ([]byte, error) {
	buf := bytes.Buffer{}
	r := bufio.NewReader(c.conn)

	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return buf.Bytes(), err
		}
		if _, err := buf.Write(line); err != nil {
			return buf.Bytes(), err
		}

		data := buf.Bytes()
		if len(data) > 2 && data[len(data)-2] == '\r' {
			return data, nil
		}
	}
	return buf.Bytes(), nil
}

func (c *Conn) ReadMessage() (*Message, error) {
	line, err := c.ReadLine()
	if err != nil {
		return nil, err
	}
	if bytes.HasSuffix(line, []byte("\r\n")) {
		line = line[:len(line)-2]
	}
	return ParseMessage(line)
}

func (c *Conn) Write(msg []byte) (int, error) {
	w := bufio.NewWriter(c.conn)
	total, err := w.Write(msg)
	if err != nil {
		return total, err
	}
	if !bytes.HasSuffix(msg, []byte("\r\n")) {
		n, err := w.WriteString("\r\n")
		total += n
		if err != nil {
			return total, err
		}
	}
	err = w.Flush()
	return total, err
}

func (c *Conn) Send(msg string) error {
    _, err := c.Write([]byte(msg))
    return err
}

func (c *Conn) Close() {
	c.conn.Close()
}
