package main

import (
	"log"
	"net"
	"time"

	"github.com/abdimk/Mort/cache"
)

func main(){
	options := SeverOptions{
		ListenAddr: ":3000",
		IsLeader: true,
	}

	go func () {
		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte("SET foo bar 2500"))
	}()

	server := NewServer(options, cache.New())
	server.Start()
}


