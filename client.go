package irc

import (
	"bufio"
	"fmt"
	"log"
	"io"
	"net/textproto"
	"strings"
	"time"
)

type Handler interface {
	Initialize(chan *Message) (err error)
	HandleMessage(*Message) (err error)
}

type Client struct {
	connection  io.ReadWriteCloser
	rw          *bufio.Reader
	handlers    []Handler
	verbose     bool
	MessageChan chan *Message
	Terminate   chan bool
}

func Connect(address string, verbose bool) (*Client, error) {
	connection, err := Dial(address)
	if err != nil {
		return nil, err
	}
	c := &Client{
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
func (c *Client) receiveLoop() {
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
func (c *Client) Serve() {
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
func (c *Client) send(message *Message) {
	fmt.Fprintf(c.connection, "%s %s :%s\r\n", message.command,
		strings.Join(message.params, " "), message.trailing)

	if c.verbose {
		log.Println(message.raw)
	}
}

// Synchronize access to connection using MessageChan channel.
func (c *Client) Send(format string, args ...interface{}) {
	raw := fmt.Sprintf(format, args...)

	if !strings.HasSuffix(raw, "\r\n") {
		raw = fmt.Sprintf("%s\r\n", raw)
	}

	message, _ := ParseLine(raw)

	c.MessageChan <- message
}

// Provide legacy interface for compatibility with echobot and managing connection sequence
// timing issues.
func (c *Client) SendImmediate(format string, args ...interface{}) {
	fmt.Fprintf(c.connection, format, args...)
	if !strings.HasSuffix(format, "\r\n") {
		fmt.Fprint(c.connection, "\r\n")
	}
}

func (c *Client) ReadMessage() (*Message, error) {
	rd := textproto.NewReader(c.rw)
	line, err := rd.ReadLine()
	if err != nil {
		return nil, err
	}
	return ParseLine(line)
}

func (c *Client) RegisterHandler(handler Handler) {
	c.handlers = append(c.handlers, handler)
}

func (c *Client) GetHandlers() (handlers []Handler) {
	return c.handlers
}

// Good news: Client is a Handler. Here we can handle details of IRC protocol
// such as server PING commands.
func (client *Client) Initialize(ch chan *Message) (err error) {
	return nil
}

func (client *Client) HandleMessage(message *Message) (err error) {
	if message.Command() == "PING" {
		client.Send("PONG %s", message.Trailing())
	}

	if client.verbose {
		log.Println(message.raw)
	}

	return nil
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
