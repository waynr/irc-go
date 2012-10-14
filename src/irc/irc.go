package irc

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"os"
)

const (
	READ_CHAN_SIZE  = 32
	WRITE_CHAN_SIZE = 32
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "irc", log.Ltime)
}

type Conn struct {
	conn     *net.TCPConn
	Received chan *Message
	ToSend   chan string
	Error    chan error
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

	c := &Conn{
		conn:     conn,
		Received: make(chan *Message, READ_CHAN_SIZE),
		ToSend:   make(chan string, WRITE_CHAN_SIZE),
		Error:    make(chan error),
	}

	// Reading task
	go func() {
		r := bufio.NewReader(conn)
		defer logger.Println("Stopping reading task")

		for {
			data, err := r.ReadString('\n')
			if err != nil {
                if err != io.EOF {
                    logger.Println("Read error: ", err)
                    c.Error <- err
                }
                c.Close()
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
					c.Error <- err
                    c.Close()
					return
				}
				conn, err := net.DialTCP("tcp", nil, addr)
				if err != nil {
					c.Error <- err
                    c.Close()
					return
				}
				c.conn = conn
			}
		}
	}()

	// Writing task
	go func() {
		w := bufio.NewWriter(conn)
        for data := range c.ToSend {
            if _, err := w.WriteString(data); err != nil {
                logger.Println("Write error: ", err)
                c.Error <- err
                return
            }
            if _, err := w.WriteString("\r\n"); err != nil {
                logger.Println("Write error: ", err)
                c.Error <- err
                return
            }
            w.Flush()
        }
	}()

	return c, nil
}

func (c *Conn) Close() {
    c.conn.Close()
    close(c.ToSend)
    close(c.Received)
}
