package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/abdimk/Mort/cache"
)

type SeverOptions struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	SeverOptions
	cache cache.Cacher
	followers map[net.Conn]struct{}
}

func NewServer(options SeverOptions, c cache.Cacher) *Server{
	return &Server{
		SeverOptions: options,
		cache: c,
		// TODO: only allocate this when we are the leader.
		followers: make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error{
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil{
		return fmt.Errorf("listen error: %s ", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	if !s.IsLeader {
		go func (){
			conn, err := net.Dial("tcp", s.LeaderAddr)
			if err != nil {
				if err == io.EOF{
					log.Println("leader closed connection")
				}else{
					log.Println("read error:", err)
				}
				return
			}
			fmt.Println("connected with leader:", s.LeaderAddr)
			s.handelConnection(conn)
		}()
	
	}

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
	defer conn.Close()


	fmt.Println("connection made:", conn.RemoteAddr())
	buf := make([]byte, 2048)


	for {
		n ,err := conn.Read(buf)

		if err != nil {
			log.Printf("conn read error: %s\n", err)
			break
		}

		data := make([]byte, n)
		copy(data, buf[:n])
		s.handleCommand(conn, data)
	}

}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte){
	msg, err := parseMessage(rawCmd)
	if err != nil{
		fmt.Println("failed to parse command", err)
		conn.Write([]byte(err.Error()))
		return
	}
	switch msg.Cmd{
		case CMDSet:
			err = s.handelSetCmd(conn, msg)
		case CMDGet:
			err = s.handelGetCmd(conn, msg)			
	}

	if err != nil {
		fmt.Println("failed to parse command:", err)
		conn.Write([]byte(err.Error()))	
	}

	
}

func (s *Server) handelGetCmd(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.key)

	if err != nil{
		return err
	}

	_,err = conn.Write(val)

	return err
}

func (s *Server) handelSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.key, msg.Value, msg.TTL); err != nil{
		return err
	}
	
	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		_, err := conn.Write(msg.ToBytes())

		if err != nil{
			fmt.Println("write to followers error:", err)
			continue
		}
	}
	return nil
}


