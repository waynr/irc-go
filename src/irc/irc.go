package irc

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"net"
)

type Conn struct {
	conn     *net.TCPConn
	Received chan *Message
	ToSend   chan string
}

func Dial(server string) (*Conn, error) {
	ipAddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, ipAddr)
	if err != nil {
		return nil, err
	}

	r := make(chan *Message, 200)
	w := make(chan string, 200)
	c := &Conn{conn: conn, Received: r, ToSend: w}

	// Reading task
	go func() {
		r := bufio.NewReader(conn)
		for {
			data, err := r.ReadString('\n')
			if err != nil {
				log.Println("Read error: ", err)
				return
			}
			msg := ParseMessage(data[0 : len(data)-2])
            c.Received <- msg
			switch {
			case bytes.Equal(msg.Command, []byte("PING")):
				c.ToSend <- "PONG" + data[4:len(data)-2]
			case bytes.Equal(msg.Command, []byte("5")):
				// RFC2812
				// Sent by the server to a user to suggest an alternative
				// server, sometimes used when the connection is refused
				// because the server is already full. Also known as RPL_SLINE
				// (AustHex), and RPL_REDIR
                addr, err := net.ResolveTCPAddr("tcp", string(msg.Trailing))
                if err != nil {
                    panic(err)
                }
				conn, err := net.DialTCP("tcp", nil, addr)
				if err != nil {
					panic(err)
				}
				c.conn = conn
			}
		}
	}()

	// Writing task
	go func() {
		w := bufio.NewWriter(conn)
		for {
			data, ok := <-c.ToSend
			if !ok {
				return
			}
			_, err := w.WriteString(data + "\r\n")
			if err != nil {
				log.Println("Write error: ", err)
			}
			w.Flush()
		}
	}()

	return c, nil
}

func (c *Conn) Close() {
}

func (c *Conn) Write(data string) error {
	c.ToSend <- data
	return nil
}

func (c *Conn) Read() (*Message, error) {
	// blocks until message is available
	data, ok := <-c.Received
	if !ok {
		return nil, errors.New("Read stream closed")
	}
	return data, nil
}
