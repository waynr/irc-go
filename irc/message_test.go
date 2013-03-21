package irc

import (
	"bytes"
	"testing"
)

func b(s string) []byte {
	return []byte(s)
}

func TestParseMessage_QUIT(t *testing.T) {
	raw := b(":Foo!bar@localhost QUIT :Bye bye!")
	m, err := ParseMessage(raw)
    if err != nil {
        t.Errorf("Message parsing error: %s", err)
    }
	if !bytes.Equal(m.Prefix, b("Foo!bar@localhost")) {
		t.Errorf("QUIT message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Params, b("")) {
		t.Errorf("QUIT message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Command, b("QUIT")) {
		t.Errorf("QUIT message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Trailing, b("Bye bye!")) {
		t.Errorf("QUIT message trailing params parsed incorrectly: %s", m.Trailing)
	}
}

func TestParseMessage_JOIN(t *testing.T) {
	raw := b(":Foo!bar@localhost JOIN #mychannel")
	m, err := ParseMessage(raw)
    if err != nil {
        t.Errorf("Message parsing error: %s", err)
    }
	if !bytes.Equal(m.Prefix, b("Foo!bar@localhost")) {
		t.Errorf("JOIN message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, b("JOIN")) {
		t.Errorf("JOIN message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, b("#mychannel")) {
		t.Errorf("JOIN message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, b("")) {
		t.Errorf("JOIN message trailing params parsed incorrectly: %s", m.Trailing)
	}
}

func TestParseMessage_MODE(t *testing.T) {
	raw := b(":Foo!bar@localhost MODE #mychannel -l")
	m, err := ParseMessage(raw)
    if err != nil {
        t.Errorf("Message parsing error: %s", err)
    }
	if !bytes.Equal(m.Prefix, b("Foo!bar@localhost")) {
		t.Errorf("MODE message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, b("MODE")) {
		t.Errorf("MODE message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, b("#mychannel -l")) {
		t.Errorf("MODE message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, b("")) {
		t.Errorf("MODE message trailing params parsed incorrectly: %s", m.Trailing)
	}
}

func TestParseMessage_PING(t *testing.T) {
	raw := b("PING :irc.localhost.localdomain")
	m, err := ParseMessage(raw)
    if err != nil {
        t.Errorf("Message parsing error: %s", err)
    }
	if !bytes.Equal(m.Prefix, b("")) {
		t.Errorf("PING message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, b("PING")) {
		t.Errorf("PING message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, b("")) {
		t.Errorf("PING message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, b("irc.localhost.localdomain")) {
		t.Errorf("PING message trailing params parsed incorrectly: %s", m.Trailing)
	}
}

func TestParseMessage_NUMERIC(t *testing.T) {
	raw := b(":wright.freenode.net 372 mynick :-")
	m, err := ParseMessage(raw)
    if err != nil {
        t.Errorf("Message parsing error: %s", err)
    }
	if !bytes.Equal(m.Prefix, b("wright.freenode.net")) {
		t.Errorf("NUMERIC message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, b("372")) {
		t.Errorf("NUMERIC message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, b("mynick")) {
		t.Errorf("NUMERIC message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, b("-")) {
		t.Errorf("NUMERIC message trailing params parsed incorrectly: %s", m.Trailing)
	}
}

func TestParseMessage_PRIVMSG(t *testing.T) {
	raw := b(":Foo!bar@s1.zz.com PRIVMSG #chan :woot")
	m, err := ParseMessage(raw)
    if err != nil {
        t.Errorf("Message parsing error: %s", err)
    }
	if !bytes.Equal(m.Prefix, b("Foo!bar@s1.zz.com")) {
		t.Errorf("PRIVMSG message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, b("PRIVMSG")) {
		t.Errorf("PRIVMSG message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, b("#chan")) {
		t.Errorf("PRIVMSG message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, b("woot")) {
		t.Errorf("PRIVMSG message trailing params parsed incorrectly: %s", m.Trailing)
	}
	if !bytes.Equal(m.Nick(), b("Foo")) {
		t.Errorf("PRIVMSG message nick parsed incorrectly: %s", m.Trailing)
	}
}
