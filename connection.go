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

// Connection represents a connection, as returned by net.Dial(), to an IRC server with
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

// Serve initiates the incoming message handling loop to pass incoming Message
// instances to all available Message handlers registered with the Connection
// object.
//
// Serve also initiates the loop which, after Serve is run, synchronizes all
// outgoing messages passed through the Connection.MessageChan channel by
// handlers.
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

// Send is used internall to put a message on the wire.
//
func (c *Connection) send(message *Message) {
	fmt.Fprintf(c.connection, "%s %s :%s\r\n", message.command,
		strings.Join(message.params, " "), message.trailing)

	if c.verbose {
		log.Println(message.raw)
	}
}

// Queue synchronizes access to the given Connection using MessageChan channel.
//
func (c *Connection) Queue(format string, args ...interface{}) {
	raw := fmt.Sprintf(format, args...)

	if !strings.HasSuffix(raw, "\r\n") {
		raw = fmt.Sprintf("%s\r\n", raw)
	}

	message, _ := ParseLine(raw)

	c.MessageChan <- message
}

// Send is intended to be used by clients to pass IRC messages to the server
// represented by Connection. This happens directly whenever the loops initiated
// by Serve are not running--that is to say, a call to Send will directly pass
// the given parameters to the underlying connection object through fmt.Fprintf
//
// After the loops initiated by Serve have begun, Send becomes and alias for
// Queue in order to synchronize access to the underlying connection object with
// whatever handlers may be registered.
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

// ReadMessage allows direct access to a bufio.Reader. It returns an irc.Message
// representing the next line in the incoming buffer.
//
func (c *Connection) ReadMessage() (*Message, error) {
	rd := textproto.NewReader(c.rw)
	line, err := rd.ReadLine()
	if err != nil {
		return nil, err
	}
	return ParseLine(line)
}

// RegisterHandler adds a Handler intended to accept incoming messages from the
// IRC server represented by the Connection.
//
func (c *Connection) RegisterHandler(handler Handler) {
	c.handlers = append(c.handlers, handler)
}

// GetHandlers returns a slice representing the Handlers currently handling
// messages from the Connection.
//
func (c *Connection) GetHandlers() (handlers []Handler) {
	return c.handlers
}

// Initialize currently does nothing, but exists as a way to allow Connection
// objects to be treated as Handler's
//
func (c *Connection) Initialize(ch chan *Message) (err error) {
	return nil
}

// HandleMessage provides a convenient way to handle details of IRC protocol
// such as server PING commands, passed in by the Connection object's incoming
// message handling loop.
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
