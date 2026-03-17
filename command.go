package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
	CMDDel Command = "DEL"
)

type MSGSet struct {
	Key   []byte
	Value []byte
	TTL   time.Duration
}

type MSGGet struct {
	Key []byte
}

type Message struct {
	Cmd   Command
	key   []byte
	Value []byte
	TTL   time.Duration
}

func (m *Message) ToBytes() []byte {
	switch m.Cmd {
	case CMDSet:
		cmd := fmt.Sprintf("%s %s %s %d", m.Cmd, m.key, m.Value, int64(m.TTL.Seconds()))
		return []byte(cmd)
	case CMDGet:
		cmd := fmt.Sprintf("%s %s", m.Cmd, m.key)
		return []byte(cmd)

	default:
		panic("unknown command")
	}
}

func parseMessage(raw []byte) (*Message, error) {
	var (
		rawStr = strings.TrimSpace(string(raw))
		parts  = strings.SplitN(rawStr, " ", 4) // Split into max 4 parts: CMD, KEY, VALUE, TTL
	)
	if len(parts) < 2 {
		return nil, errors.New("invalid protocol format")
	}

	msg := &Message{
		Cmd: Command(parts[0]),
		key: []byte(parts[1]),
	}

	if msg.Cmd == CMDSet {
		if len(parts) < 4 {
			return nil, errors.New("invalid SET command")
		}

		msg.Value = []byte(parts[2])

		ttl, err := time.ParseDuration(parts[3])

		if err != nil {
			sec, err2 := strconv.Atoi(parts[3])
			if err2 != nil {
				return nil, err
			}
			ttl = time.Duration(sec) * time.Second
		}

		msg.TTL = ttl
	}

	return msg, nil

}