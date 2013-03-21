package irc

import (
	"bufio"
	"bytes"
	"net"
	"testing"
)

func echoIRCServer(t *testing.T) (net.Listener, chan error) {
	c, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Errorf("Cannot create test server: %s", err)
	}
	quit := make(chan error, 2)

	go func() {
		for {
			cli, err := c.Accept()
			if err != nil {
				quit <- err
				return
			}
			go handleClient(cli, quit)
		}
	}()

	errchan := make(chan error, 1)

	go func() {
		err := <-quit
		quit <- err
		if err != nil {
			errchan <- err
		}
		c.Close()
	}()

	return c, errchan
}

func handleClient(c net.Conn, quit chan error) {
	r := bufio.NewReader(c)
	buf := bytes.Buffer{}
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			quit <- err
			return
		}
		buf.Write(line)
		if line[len(line)-2] == '\r' {
		}
		if _, err := c.Write(buf.Bytes()); err != nil {
			quit <- err
			return
		}
		buf.Reset()
	}
}
