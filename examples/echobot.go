package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/husio/go-irc"
)

var (
	address = flag.String("address", "irc.freenode.net:6667", "IRC server address")
	nick    = flag.String("nick", "gobot", "User nick")
	name    = flag.String("name", "GoBot", "User name")
	verbose = flag.Bool("verbose", false, "Print all messages to stdout")
)

func main() {
	flag.Parse()
	c, err := irc.Connect(*address, *verbose)
	if err != nil {
		log.Fatalln(err)
	}

	c.SendImmediate("NICK %s", *nick)
	c.SendImmediate("USER %s * * :...", *name)
	time.Sleep(time.Millisecond * 10)

	for _, name := range flag.Args() {
		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}
		c.SendImmediate("JOIN %s", name)
	}

	// read data from stdin and send it through the wire
	go func() {
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
			c.SendImmediate(line)
		}
	}()

	fmt.Print(`

For IRC protocol description, read rfc1459: https://tools.ietf.org/html/rfc1459
Some basics:

	JOIN    <channel>{,<channel>} [<key>{,<key>}]
	PRIVMSG <receiver>{,<receiver>} <text to be sent>
	PART    <channel>{,<channel>}

`)

	// handle incomming messages
	for {
		message, err := c.ReadMessage()
		if err != nil {
			log.Fatalln(err)
			return
		}
		if message.Command() == "PING" {
			c.SendImmediate("PONG %s", message.Trailing())
		}

		if message.Command() == "PRIVMSG" {
			if strings.HasPrefix(message.Trailing(), *nick) {
				text := message.Trailing()[len(*nick):]
				c.SendImmediate("PRIVMSG %s :echo! \"%s\"", message.Params()[0], text)
			}
		}

		if *verbose {
			log.Println(message)
		}
	}
}
