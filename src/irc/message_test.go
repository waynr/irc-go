package irc

import (
   "bytes"
    "testing"
)

func TestParseMessage_QUIT(t *testing.T) {
	raw := ":Foo!bar@localhost QUIT :Bye bye!"
	m := ParseMessage(raw)
	if !bytes.Equal(m.Prefix, []byte("Foo!bar@localhost")) {
		t.Errorf("QUIT message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Params, []byte("")) {
		t.Errorf("QUIT message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Command, []byte("QUIT")) {
		t.Errorf("QUIT message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Trailing, []byte("Bye bye!")) {
		t.Errorf("QUIT message trailing params parsed incorrectly: %s",
			m.Trailing)
	}
}

func TestParseMessage_JOIN(t *testing.T) {
	raw := ":Foo!bar@localhost JOIN #mychannel"
	m := ParseMessage(raw)
	if !bytes.Equal(m.Prefix, []byte("Foo!bar@localhost")) {
		t.Errorf("JOIN message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, []byte("JOIN")) {
		t.Errorf("JOIN message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, []byte("#mychannel")) {
		t.Errorf("JOIN message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, []byte("")) {
		t.Errorf("JOIN message trailing params parsed incorrectly: %s",
			m.Trailing)
	}
}

func TestParseMessage_MODE(t *testing.T) {
	raw := ":Foo!bar@localhost MODE #mychannel -l"
	m := ParseMessage(raw)
	if !bytes.Equal(m.Prefix, []byte("Foo!bar@localhost")) {
		t.Errorf("MODE message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, []byte("MODE")) {
		t.Errorf("MODE message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, []byte("#mychannel -l")) {
		t.Errorf("MODE message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, []byte("")) {
		t.Errorf("MODE message trailing params parsed incorrectly: %s",
			m.Trailing)
	}
}

func TestParseMessage_PING(t *testing.T) {
	raw := "PING :irc.localhost.localdomain"
	m := ParseMessage(raw)
	if !bytes.Equal(m.Prefix, []byte("")) {
		t.Errorf("PING message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, []byte("PING")) {
		t.Errorf("PING message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, []byte("")) {
		t.Errorf("PING message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, []byte("irc.localhost.localdomain")) {
		t.Errorf("PING message trailing params parsed incorrectly: %s",
			m.Trailing)
	}
}

func TestParseMessage_NUMERIC(t *testing.T) {
	raw := ":wright.freenode.net 372 mynick :-"
	m := ParseMessage(raw)
	if !bytes.Equal(m.Prefix, []byte("wright.freenode.net")) {
		t.Errorf("NUMERIC message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, []byte("372")) {
		t.Errorf("NUMERIC message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, []byte("mynick")) {
		t.Errorf("NUMERIC message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, []byte("-")) {
		t.Errorf("NUMERIC message trailing params parsed incorrectly: %s",
			m.Trailing)
	}
}

func TestParseMessage_PRIVMSG(t *testing.T) {
    raw := ":Foo!bar@s1.zz.com PRIVMSG #chan :woot"
	m := ParseMessage(raw)
	if !bytes.Equal(m.Prefix, []byte("Foo!bar@s1.zz.com")) {
		t.Errorf("PRIVMSG message prefix parsed incorrectly: %s", m.Prefix)
	}
	if !bytes.Equal(m.Command, []byte("PRIVMSG")) {
		t.Errorf("PRIVMSG message type parsed incorrectly: %d", m.Command)
	}
	if !bytes.Equal(m.Params, []byte("#chan")) {
		t.Errorf("PRIVMSG message params parsed incorrectly: %s", m.Params)
	}
	if !bytes.Equal(m.Trailing, []byte("woot")) {
		t.Errorf("PRIVMSG message trailing params parsed incorrectly: %s",
			m.Trailing)
	}
	if !bytes.Equal(m.Nick(), []byte("Foo")) {
		t.Errorf("PRIVMSG message nick parsed incorrectly: %s",
			m.Trailing)
	}
}
