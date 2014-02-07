package irc

import (
	"bufio"
	"fmt"
	"log"
	"io"
	"net"
	"net/textproto"
	"strings"
	"time"
)

// Represents a connection, as returned by net.Dial(), to an IRC server with
// methods for synchronizing access.
//
type Connection struct {
	connection  io.ReadWriteCloser
	rw          *bufio.Reader
	handlers    []Handler
	verbose, serving     bool

	// MessageChan is passed to handlers to provide a way for handlers to
	// synchronize write access to the underlying IRC connection.
	MessageChan chan *Message

	// Terminate is used to signal the underlying IRC connection is ready to be
	// shut down.
	Terminate   chan bool
}

// Connect to the given address and return a connected Connection object.
//
func Connect(address string, verbose bool) (*Connection, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	c := &Connection{
		connection:  connection,
		rw:          bufio.NewReader(connection),
		handlers:    make([]Handler, 0, 4),
		verbose:     verbose,
		MessageChan: make(chan *Message, 100),
		Terminate:   make(chan bool),
	}

	c.RegisterHandler(c)
	return c, nil
}

// Manage incoming messages from the server by dispatching them to registered
// handlers.
//
func (c *Connection) receiveLoop() {
	for {
		message, err := c.ReadMessage()
		if err != nil {
			log.Fatalln(err)
			return
		}

		for _, handler := range c.handlers {
			handler.HandleMessage(message)
		}

	}
}

// Manage MessageChan and Terminate channels in a forever loop.
//
func (c *Connection) Serve() {
	c.serving = true
	defer func() {c.serving = false}()
	go c.receiveLoop()

	for {
		select {
		case message, ok := <-c.MessageChan:
			if ok {
				c.send(message)
			} else {
				break
			}
		case <- c.Terminate:
			close(c.MessageChan)
		default:
		}
		time.Sleep(time.Millisecond *10)
	}
}

// Put a message on the wire.
//
func (c *Connection) send(message *Message) {
	fmt.Fprintf(c.connection, "%s %s :%s\r\n", message.command,
		strings.Join(message.params, " "), message.trailing)

	if c.verbose {
		log.Println(message.raw)
	}
}

// Synchronize access to the Connection using MessageChan channel.
//
func (c *Connection) Queue(format string, args ...interface{}) {
	raw := fmt.Sprintf(format, args...)

	if !strings.HasSuffix(raw, "\r\n") {
		raw = fmt.Sprintf("%s\r\n", raw)
	}

	message, _ := ParseLine(raw)

	c.MessageChan <- message
}

// While the Connection is in a "serving" state, that is to say messages are
// being synchronously read from its MessageChan, Send() becomes an alias for the
// Queue() method
//
// While the Connection is static, that is to say, it is not synchronizing write
// access to the MessageChan, it writes directly to the underlying connection
// object.
//
func (c *Connection) Send(format string, args ...interface{}) {
	if c.serving {
		c.Queue(format, args...)
	} else {
		fmt.Fprintf(c.connection, format, args...)
		if !strings.HasSuffix(format, "\r\n") {
			fmt.Fprint(c.connection, "\r\n")
		}
	}
}

// When the Connection is not serving, ReadMessage allows direct access to a
// bufio.Reader. It returns an irc.Message representing the next line in the
// incoming buffer.
//
func (c *Connection) ReadMessage() (*Message, error) {
	rd := textproto.NewReader(c.rw)
	line, err := rd.ReadLine()
	if err != nil {
		return nil, err
	}
	return ParseLine(line)
}

// Adds a Handler intended to accept incoming messages from the IRC server
// represented by the irc.Connection.
//
func (c *Connection) RegisterHandler(handler Handler) {
	c.handlers = append(c.handlers, handler)
}

// Returns the slice representing the Handlers currently handling messages from
// the irc.Connection.
//
func (c *Connection) GetHandlers() (handlers []Handler) {
	return c.handlers
}

func (c *Connection) Initialize(ch chan *Message) (err error) {
	return nil
}

// Connection is a Handler. This provides a convenient place way to handle
// details of IRC protocol such as server PING commands.
//
func (c *Connection) HandleMessage(message *Message) (err error) {
	if message.Command() == "PING" {
		c.Queue("PONG %s", message.Trailing())
	}

	if c.verbose {
		log.Println(message.raw)
	}

	return nil
}

/*
	switch msg.Command {
	case "PING":
		c.Queue("PONG %s\r\n", msg.Trailing)
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
