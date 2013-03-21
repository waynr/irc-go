package irc

import (
	"bytes"
	"testing"
)

func TestReadWrite(t *testing.T) {
	srv, _ := echoIRCServer(t)

	cli, err := Dial(srv.Addr().String())
	if err != nil {
		t.Fatalf("Cannot connect to server: %s", err)
	}
	defer cli.Close()

	cli.Write([]byte("foo\nbaz\rbar"))
	resp, err := cli.ReadLine()
	if err != nil {
		t.Fatal("ReadLine error", err)
	}
	if !bytes.Equal(resp, []byte("foo\nbaz\rbar\r\n")) {
		t.Fatalf("Expected 'foo\\nbaz\\rbar\\r\\n', got '%s'", resp)
	}
}

func TestLostConnection(t *testing.T) {
	srv, _ := echoIRCServer(t)

	cli, err := Dial(srv.Addr().String())
	if err != nil {
		t.Fatalf("Cannot connect to server: %s", err)
	}

	// XXX - we want to test with timeout here, but it just takes ages!
	srv.Close()
	cli.Close()

	if _, err := cli.ReadLine(); err == nil {
		t.Fatal("Reading from broken connection succeded?")
	}
	if _, err := cli.Write([]byte("foo")); err == nil {
		t.Fatal("Writing to broken connection succeded?")
	}
}

func TestMessage(t *testing.T) {
	srv, _ := echoIRCServer(t)

	cli, err := Dial(srv.Addr().String())
	if err != nil {
		t.Fatalf("Cannot connect to server: %s", err)
	}
	defer cli.Close()

    line := []byte(":Foo!bar@s1.zz.com PRIVMSG #chan :woot\r\n")
	cli.Write(line)
	m, err := cli.ReadMessage()
	if err != nil {
		t.Fatal("ReadMessage error", err)
	}
	if !bytes.Equal(m.raw, line) {
		t.Fatalf("Expected '%s', got '%s'", line, m)
	}
}
