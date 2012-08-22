package irc

import (
	"bytes"
)

type Message struct {
	raw      string
	Prefix   []byte
	Command     []byte
	Params   []byte
	Trailing []byte
}

func ParseMessage(message string) *Message {
	// :<prefix> <command> <params> :<trailing>
	msg := []byte(message)
	var prefix, params, trailing []byte

	if bytes.HasPrefix(msg, []byte(":")) {
		index := bytes.Index(msg, []byte(" "))
		prefix = msg[1:index]
		msg = msg[index+1:]
	}

	cmdEndIndex := bytes.Index(msg, []byte(" "))
	command := msg[:cmdEndIndex]
	msg = msg[cmdEndIndex+1:]
	trailingStartIndex := bytes.Index(msg, []byte(":"))
	if trailingStartIndex < 0 {
		params = msg
	} else {
		params = bytes.TrimRight(msg[:trailingStartIndex], " ")
		trailing = msg[trailingStartIndex+1:]
	}

	return &Message{
		raw:      message,
		Prefix:   prefix,
		Command:     command,
		Params:   params,
		Trailing: trailing,
	}
}

func (msg *Message) Nick() []byte {
	expIndex := bytes.Index(msg.Prefix, []byte("!"))
	if expIndex == -1 {
		return nil
	}
	return msg.Prefix[:expIndex]
}

func (msg *Message) String() string {
    return msg.raw
}
