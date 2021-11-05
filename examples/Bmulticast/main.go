package main

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"MulticastSDCCProject/pkg/pool"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	//var port uint
	port := flag.Uint("port", 8090, "server port number")
	group := flag.String("groupPort", "8090,8091,8092", "defining group port")

	flag.Parse()

	var err error
	go func() {

		err = server.RunServer(*port)
	}()
	if err != nil {
		log.Println("Error in connection node")
		return
	}

	var localErr error

	var groupArray []string
	var connections []*client.Client

	groupArray = strings.Split(*group, ",")

	connections = make([]*client.Client, len(groupArray))

	for i := range groupArray {
		connections[i] = client.Connect("localhost:" + groupArray[i])
	}

	respChannel := make(chan []byte, 1)
	pool.Pool.InitThreadPool(connections, 5, util.BMULTICAST, respChannel, *port)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Insert message: ")
	for scanner.Scan() {

		m := &rpc.Packet{Message: scanner.Bytes()}
		pool.Pool.Message <- m

	}

	if scanner.Err() != nil {
		log.Println("Error from stdin")
	}

	if localErr != nil {
		log.Println("Error in sending to node")
	}

}
