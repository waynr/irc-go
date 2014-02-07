package irc

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidMessage is used by irc.ParseLine to indicate that the given
	// "raw" string does not conform to RFC1459.
	ErrInvalidMessage = errors.New("invalid message format")
	// ErrUnknownCommand by irc.ParseLine to indicate error condition in which no command was
	// found in the given "raw" message.
	ErrUnknownCommand = errors.New("unknown command")
)

// Handler interface provides a means for users to implement "handlers" for
// irc.Message objects. This interface is used by irc.Connection to dispatch
// messages to custom handlers. 
//
// The user may choose to do as she wishes with these messages or to respond in
// her own time to the underlying Server by queuing messages through the
// chan *Message passed in through the Initialize method.
//
type Handler interface {
	// Used by irc.Connection to provide access to its chan *Message
	Initialize(chan *Message) (err error)

	// Used by irc.Connection to pass in messages from the server.
	HandleMessage(*Message) (err error)
}

// Message is a struct representing a go message according to
// http://tools.ietf.org/html/rfc1459.html#section-2.3.1
//
type Message struct {
	raw      string
	prefix   string
	command  string
	params   []string
	trailing string
}

// ParseLine is used to create irc.Message objects out of raw strings.
//
// When sending messages, <prefix> is automatically determined by the receiving
// end of the connection based on available TCP/IP and DNS connection
// information.  ParseLine takes this into account by making the ":<prefix> "
// optional.
//
// When used to construct messages, the raw string must follow the guidelines in
// RFC1459: http://tools.ietf.org/html/rfc1459
//
// Please note that it is NOT the user's responsibility to deal with
// prefixes--this is determined upon arrival at the IRC server or client.
// Therefore, for the purpose of this library please use the following format
// for "raw" strings, where <command> is a valied ASCII command as described in
// the RFC1459.
//
// <command>[ <params> [:<trailing>]]
//
func ParseLine(raw string) (*Message, error) {
	raw = strings.TrimSpace(raw)
	m := &Message{raw: raw}
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

// Prefix provides read-only public access to the internal IRC "prefix".
//
func (m *Message) Prefix() string {
	return m.prefix
}

// Command provides read-only public access to the internal IRC "command".
//
func (m *Message) Command() string {
	return m.command
}

// Params provides read-only public access to the internal IRC "params".
//
func (m *Message) Params() []string {
	return m.params
}

// Trailing provides read-only public access to the internal IRC "trailing".
//
func (m *Message) Trailing() string {
	return m.trailing
}

// String provides read-only public access to the internal IRC raw representation of the
// message.
//
func (m *Message) String() string {
	return m.raw
}
