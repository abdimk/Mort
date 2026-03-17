package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/abdimk/Mort/cache"
)

type SeverOptions struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	SeverOptions
	cache     cache.Cacher
	followers map[net.Conn]struct{}
}

func NewServer(options SeverOptions, c cache.Cacher) *Server {
	return &Server{
		SeverOptions: options,
		cache:        c,
		// TODO: only allocate this when we are the leader.
		followers: make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return fmt.Errorf("listen error: %s ", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	if !s.IsLeader {
		go func() {
			conn, err := net.Dial("tcp", s.LeaderAddr)
			if err != nil {
				log.Printf("failed to connect to leader %s: %s", s.LeaderAddr, err)
				return
			}
			log.Printf("connected to leader: %s", s.LeaderAddr)
			s.handleFollowerConnection(conn)
		}()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n ", err)
			continue
		}

		if s.IsLeader {
			go s.handleLeaderConnection(conn)
		}
	}
}

// handleFollowerConnection handles the connection when running as a follower
// It only listens and logs SET/GET operations from the leader, never responds
func (s *Server) handleFollowerConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("follower listening from leader: %s", conn.RemoteAddr())
	buf := make([]byte, 2048)

	for {
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				log.Printf("leader connection closed: %s", conn.RemoteAddr())
				return
			}

			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				continue // keep listening
			}

			log.Printf("follower connection error: %s", err)
			return
		}

		data := make([]byte, n)
		copy(data, buf[:n])
		s.logCommandFromLeader(data)
	}
}

// logCommandFromLeader logs SET/GET operations received from the leader
func (s *Server) logCommandFromLeader(rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		log.Printf("failed to parse command from leader: %s", err)
		return
	}

	switch msg.Cmd {
	case CMDSet:
		log.Printf("[FOLLOWER] SET key=%s value=%s ttl=%s\n", string(msg.key), string(msg.Value), msg.TTL)
	case CMDGet:
		val, err := s.cache.Get(msg.key)
		if err != nil {
			log.Printf("[FOLLOWER] GET key=%s (not found)\n", string(msg.key))
		} else {
			log.Printf("[FOLLOWER] GET key=%s = %s\n", string(msg.key), string(val))
		}
	}
}

// handleLeaderConnection handles connections when running as the leader
func (s *Server) handleLeaderConnection(conn net.Conn) {
	defer func() {
		delete(s.followers, conn)
		conn.Close()
	}()

	s.followers[conn] = struct{}{}

	log.Printf("connection made: %s", conn.RemoteAddr())
	buf := make([]byte, 2048)

	// Read once, handle command, then disconnect (echo | nc style)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)

	if err != nil {
		if err == io.EOF {
			log.Printf("connection closed: %s", conn.RemoteAddr())
			return
		}

		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			log.Printf("connection timeout: %s", conn.RemoteAddr())
			return
		}

		log.Printf("connection error: %s", err)
		return
	}

	data := make([]byte, n)
	copy(data, buf[:n])
	s.handleCommand(conn, data)
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		log.Printf("failed to parse command from %s: %s", conn.RemoteAddr(), err)
		conn.Write([]byte(err.Error()))
		return
	}
	switch msg.Cmd {
	case CMDSet:
		err = s.handelSetCmd(conn, msg)
	case CMDGet:
		err = s.handelGetCmd(conn, msg)
	}

	if err != nil {
		log.Printf("command error from %s: %s", conn.RemoteAddr(), err)
		conn.Write([]byte(err.Error()))
	}
}

func (s *Server) handelGetCmd(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.key)

	if err != nil {
		return err
	}

	_, err = conn.Write(val)

	return err
}

func (s *Server) handelSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		_, err := conn.Write(msg.ToBytes())

		if err != nil{
			log.Printf("failed to forward to follower %s: %s", conn.RemoteAddr(), err)
			delete(s.followers, conn)
			conn.Close()
			continue
		}
	}
	return nil
}


