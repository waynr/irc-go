package irc

import (
	"errors"
	"strings"
)

var (
	ErrInvalidmessage = errors.New("invalid message format")
	ErrUnknownCommand = errors.New("unknown command")
)

// IRC message format:
//
// :<prefix> <command> <params> :<trailing>
type Message interface {
	Prefix() string
	Command() string
	Params() []string
	Trailing() string
}

type message struct {
	raw      string
	prefix   string
	command  string
	params   []string
	trailing string
}

func ParseLine(raw string) (*message, error) {
	raw = strings.TrimSpace(raw)
	m := &message{raw: raw}
	if raw[0] == ':' {
		chunks := strings.SplitN(raw, " ", 2)
		m.prefix = chunks[0][1:]
		raw = chunks[1]
	}
	chunks := strings.SplitN(raw, " ", 2)
	m.command = chunks[0]
	raw = chunks[1]
	if m.command == "" {
		return nil, ErrUnknownCommand
	}

	if raw[0] != ':' {
		chunks := strings.SplitN(raw, " :", 2)
		m.params = strings.Split(chunks[0], " ")
		if len(chunks) == 2 {
			raw = chunks[1]
		} else {
			raw = ""
		}
	}

	if len(raw) > 0 {
		if raw[0] == ':' {
			raw = raw[1:]
		}
		m.trailing = raw
	}
	return m, nil
}

func (m *message) Prefix() string {
	return m.prefix
}

func (m *message) Command() string {
	return m.command
}

func (m *message) Params() []string {
	return m.params
}

func (m *message) Trailing() string {
	return m.trailing
}

func (m *message) String() string {
	return m.raw
}
