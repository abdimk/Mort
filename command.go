package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type MSGSet struct {
	Key []byte
	Value []byte
	TTL time.Duration
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

func parseMessage(raw []byte)(*Message, error){
	var (
		rawStr = string(raw)
		parts = strings.Split(rawStr, " ")
	)
	if len(parts) < 2 {
		// respond
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
		ttl, err := strconv.Atoi(parts[3])	

		if err != nil {
			return nil, err
		}

		msg.TTL = time.Duration(ttl) 
	}

	return msg, nil

}