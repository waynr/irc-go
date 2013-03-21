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

	addr := fmt.Sprintf("%s:%d", *host, *port)
	c, err := irc.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := fmt.Fprintf(c, "NICK %s", *nick); err != nil {
		log.Fatal(err)
	}
	if _, err := fmt.Fprintf(c, "USER bot * * :..."); err != nil {
		log.Fatal(err)
	}

	for _, name := range flag.Args() {
		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}
		fmt.Printf("[join] %s\n", name)
		fmt.Fprintf(c, "JOIN %s", name)
	}

	mynick := []byte(*nick)
	privcmd := []byte("PRIVMSG")
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("ReadMessage error: %s", err)
		}
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
