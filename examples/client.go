package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"../irc"
)

var host *string = flag.String("host", "irc.freenode.net", "IRC server address")
var port *int = flag.Int("port", 6667, "IRC server port")
var nick *string = flag.String("nick", "go-irc-client", "Nickname")

var help = `
********************************************************************************

quit                               - close the client


JOIN #<name> 					   - join channel
PRIVMSG #<channel name> :<message> - send message to given channel


More info: http://tools.ietf.org/html/rfc1459

********************************************************************************
`

func main() {
	flag.Parse()

	c, err := irc.Connect(*host, *port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Printf("\n** For more information type `help` **\n\n")

	if err := c.Send("NICK %s", *nick); err != nil {
		log.Fatal(err)
	}
	if err := c.Send("USER bot * * :..."); err != nil {
		log.Fatal(err)
	}

	// reader
	go func() {
		for {
			select {
			case err := <-c.Error:
				log.Fatalf("IRC client error: %s", err)
			case msg := <-c.Received:
				fmt.Printf("[message] %s\n", msg)
			}
		}
	}()

	// writer
	in := bufio.NewReader(os.Stdin)
	for {
		data, err := in.ReadString('\n')
		if err != nil {
			log.Fatal("Client write error: %s", err)
		}
		data = strings.TrimSpace(data)
		switch data {
		case "help":
			fmt.Println(help)
		case "quit":
			return
		default:
			err := c.Send(data)
			if err != nil {
				log.Fatal("Sending message error: %s", err)
			}
		}
	}
}
