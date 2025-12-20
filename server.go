package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/abdimk/Mort/cache"
)

type SeverOptions struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	SeverOptions
	cache cache.Cacher
}

func NewServer(options SeverOptions, c cache.Cacher) *Server{
	return &Server{
		SeverOptions: options,
		cache: c,
	}
}

func (s *Server) Start() error{
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil{
		return fmt.Errorf("listen error: %s ", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)


	for {
		conn, err := ln.Accept()
		if err != nil{
			log.Printf("accept error: %s\n ", err)
			continue
		}

		go s.handelConnection(conn)
	}
}

func (s *Server) handelConnection(conn net.Conn){
	defer func(){
		conn.Close()
	}()

	buf := make([]byte, 2048)


	for {
		n ,err := conn.Read(buf)

		if err != nil {
			log.Printf("conn read error: %s\n", err)
			break
		}

		data := make([]byte, n)
		copy(data, buf[:n])
		go s.handleCommand(conn, data)
	}

}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte){
	msg, err := parseMessage(rawCmd)
	if err != nil{
		fmt.Println("failed to parse command", err)
		return
	}

	switch msg.Cmd{
		case CMDSet:
			if err := s.handelSetCommand(conn, msg); err != nil {
				// respond
				return
			}
	}
	
}

func (s *Server) handelSetCommand(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.key, msg.Value, msg.TTL); err != nil{
		return err
	}
	
	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}


