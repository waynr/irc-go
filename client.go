package irc

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strings"
)

type Handler interface {
}

type Client interface {
	Send(format string, args ...interface{})
	RegisterHandler(handler Handler)
}

type client struct {
	connection io.ReadWriteCloser
	rw         *bufio.Reader
	wr         *bufio.Writer
	handlers   []Handler
}

func Connect(address string) (*client, error) {
	connection, err := Dial(address)
	if err != nil {
		return nil, err
	}
	c := &client{
		connection: connection,
		rw:         bufio.NewReader(connection),
		wr:         bufio.NewWriter(connection),
		handlers:   make([]Handler, 0, 4),
	}
	return c, nil
}

func (c *client) Send(format string, args ...interface{}) {
	fmt.Fprintf(c.connection, format, args...)
	if !strings.HasSuffix(format, "\r\n") {
		fmt.Fprint(c.connection, "\r\n")
	}
}

func (c *client) ReadMessage() (*message, error) {
	rd := textproto.NewReader(c.rw)
	line, err := rd.ReadLine()
	if err != nil {
		return nil, err
	}
	return ParseLine(line)
}

/*
	switch msg.Command {
	case "PING":
		c.Send("PONG %s\r\n", msg.Trailing)
	case "JOIN":
		for _, name := range strings.Split(msg.Params[0], ",") {
			c.channels[name] = name
		}
	case "NICK":
		c.nick = msg.Trailing
	case "PART":
		for _, name := range strings.Split(msg.Params[0], ",") {
			delete(c.channels, name)
		}
	case "KICK":
		if msg.Params[1] == c.nick {
			delete(c.channels, msg.Params[0])
		}
	case "001":
		// "Welcome to the Internet Relay Network <nick>!<user>@<host>"
		rx := regexp.MustCompile(`(\w+)(\!|$)`)
		nick := rx.FindString(line)
		if nick[len(nick) - 1] == '!' {
			nick = nick[:len(nick) - 1]
		}
		c.nick = nick
	}
	return msg, nil
*/
