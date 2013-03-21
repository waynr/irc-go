package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"../irc"
)

var host *string = flag.String("host", "irc.freenode.net", "IRC server address")
var port *int = flag.Int("port", 6667, "IRC server port")
var nick *string = flag.String("nick", "go-echobot", "Nickname")

const usage = `Usage:

%s [<flags>] <channel name1> [<channel name 2>], ...

`

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf(usage, os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	c, err := irc.Connect(*host, *port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if err := c.Send("NICK %s", *nick); err != nil {
		log.Fatal(err)
	}
	if err := c.Send("USER bot * * :..."); err != nil {
		log.Fatal(err)
	}

	for _, name := range flag.Args() {
		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}
		fmt.Printf("[join] %s\n", name)
		if err := c.Send("JOIN %s", name); err != nil {
			log.Fatal(err)
		}
	}

	mynick := []byte(*nick)
	privcmd := []byte("PRIVMSG")
	for {
		select {
		case err := <-c.Error:
			log.Fatalf("IRC client error: %s", err)
		case msg := <-c.Received:
			fmt.Printf("[message] %s\n", msg)

			if bytes.Equal(msg.Command, privcmd) && bytes.HasPrefix(msg.Trailing, mynick) && len(msg.Trailing) > len(mynick)+2 {

				text := msg.Trailing[len(mynick)+2:]
				if bytes.Equal(text, []byte("foo")) {
					text = []byte("bar ;)")
				}

				resp := fmt.Sprintf("PRIVMSG %s :%s: %s", msg.Params, msg.Nick(), text)
				fmt.Printf("[echo response] %s\n", resp)
				if err := c.Send(resp); err != nil {
					log.Fatal("Sending echo error: %s", err)
				}
			}
		}
	}
}
