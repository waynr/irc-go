package irc

import (
	"bytes"
	"fmt"
)

type Message struct {
	raw      []byte
	Prefix   []byte
	Command  []byte
	Params   []byte
	Trailing []byte
}

func ParseMessage(message []byte) (*Message, error) {
	var err error

	defer func () {
		x := recover()
		if x != nil {
			err = fmt.Errorf("Parsing error: %v", x)
		}
	}()
	msg := parse(message)
	return msg, err
}

func parse(message []byte) *Message {
	// :<prefix> <command> <params> :<trailing>
	msg := message
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
		Command:  command,
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
	return string(msg.raw)
}
