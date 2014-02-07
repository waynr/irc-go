package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/waynr/irc-go"
)

var (
	address = flag.String("address", "irc.freenode.net:6667", "IRC server address")
	nick = flag.String("nick", "handlerBot", "User nick")
	name = flag.String("name", "hbot", "User name")
	verbose = flag.Bool("verbose", false, "Print all messages to stdout.")
)

func main() {
	flag.Parse()

	c, err := irc.Connect(*address, *verbose)
	if err != nil {
		log.Fatalln(err)
	}

	c.Send("NICK %s", *nick)
	c.Send("USER %s * * :...", *name)

	time.Sleep(time.Millisecond *10)

	for _, name := range flag.Args() {
		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}
		c.Queue("JOIN %s", name)
	}

	// begin handling Stdin
	go handleStdin(c)

	// register echo handler
	c.RegisterHandler(&echoHandler{})

	// pass channels to registered handlers for handling
	for _, handler := range c.GetHandlers() {
		handler.Initialize(c.MessageChan)
	}

	c.Serve()
}

// takes input from terminal and queues to be sent to server
func handleStdin(c *irc.Connection) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		c.Queue(line)
	}
}

// Example "echo" handler
type echoHandler struct {
	MessageChan chan *irc.Message
}

func (echo *echoHandler) Initialize(ch chan *irc.Message) (err error) {
	echo.MessageChan = ch
	return nil
}

func (echo *echoHandler) HandleMessage(message *irc.Message) (err error) {
	if message.Command() == "PRIVMSG" {
		if strings.HasPrefix(message.Trailing(), *nick) {
			chunks := strings.SplitN(message.Trailing(), " ", 2)
			raw := fmt.Sprintf("PRIVMSG %s :%s", message.Params()[0], chunks[1])
			newMessage, _ := irc.ParseLine(raw)
			echo.MessageChan <- newMessage
		}
	}
	return nil
}
