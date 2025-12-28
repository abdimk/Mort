package main

import (
	"flag"

	"github.com/abdimk/Mort/cache"

)
func main(){

	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	
	flag.Parse()

	options := SeverOptions{
		ListenAddr: *listenAddr,
		IsLeader: len(*leaderAddr)==0,
		LeaderAddr: *leaderAddr,
	}

	// conn, err := net.Dial("tcp", ":3000")
	// if err != nil {
	// 	fmt.Println("error while tring to dile the leader", err.Error())
	// }

	// conn.Write([]byte("SET Foo Bar 400000"))

	// buf := make([]byte, 2048)

	// n, err := conn.Read(buf)

	// if err != nil {
	// 	fmt.Println("unable tp read the bytes", err.Error())
	// }

	// fmt.Printf("Msg: %v", string(buf[:n]))

	// comd := []byte("SET Foo Bar 40000")
	// k,err := parseMessage(comd)
	// if err != nil{
	// 	fmt.Println(err)
	// }

	// fmt.Println(reflect.TypeOf(k.TTL))


	

	server := NewServer(options, cache.New())
	server.Start()
}

 
