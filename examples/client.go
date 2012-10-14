package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"irc"
	"os"
	"strings"
)

var server *string = flag.String("server", "irc.freenode.net", "IRC server address")
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

	addr := fmt.Sprintf("%s:%v", *server, *port)
	c, err := irc.Dial(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n** For more information type `help` **\n\n")

	defer c.Close()

	quit := make(chan bool)

	c.ToSend <- "NICK " + *nick
	c.ToSend <- "USER bot * * :..."

	// irc messages reader
	go func() {
		for {
			select {
			case err := <-c.Error:
				fmt.Println("client read error", err)
				quit <- true
				return
			case msg := <-c.Received:
				if bytes.Equal(msg.Command, []byte("PRIVMSG")) {
					fmt.Printf("%s:: %s %s -> %s\n",
						msg.Command, msg.Params, msg.Prefix, msg.Trailing)
				} else {
					fmt.Println("> ", msg.String())
				}
			}
		}
	}()

	// user input reader
	go func() {
		in := bufio.NewReader(os.Stdin)
		for {
			data, err := in.ReadString('\n')
			if err != nil {
				fmt.Sprintf("client write error: %s", err)
				return
			}
			data = strings.TrimSpace(data)
			switch data {
			case "help":
				fmt.Println(help)
			case "quit":
				quit <- true
			default:
				c.ToSend <- data
			}
		}
	}()

	<-quit
}
