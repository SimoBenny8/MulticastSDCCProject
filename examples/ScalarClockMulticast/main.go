package main

import (
	"MulticastSDCCProject/pkg/MulticastScalarClock/impl"
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"MulticastSDCCProject/pkg/pool"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//var port uint
	port := flag.Uint("port", 8090, "server port number")
	group := flag.String("groupPort", "8090,8091,8092", "defining group port")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	var err error
	var localErr error
	var groupArray []string
	var connections []*client.Client
	var myConn *client.Client

	go func() {

		err = server.RunServer(*port, server.Register)
	}()
	if err != nil {
		log.Println("Error in connection node")
		return
	}

	groupArray = strings.Split(*group, ",")

	connections = make([]*client.Client, len(groupArray))

	for i := range groupArray {
		if groupArray[i] == strconv.Itoa(int(*port)) {
			connections[i] = client.Connect("localhost:" + groupArray[i])
			myConn = connections[i]
		} else {
			connections[i] = client.Connect("localhost:" + groupArray[i])
		}

	}

	pool.Pool.InitThreadPool(connections, 5, util.SCMULTICAST, nil, *port)
	go impl.Receive(connections, *port)
	go impl.Deliver(myConn, len(connections))
	impl.SetConnections(connections)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Insert message: ")
	for scanner.Scan() {

		m := &rpc.Packet{Message: scanner.Bytes()}
		pool.Pool.Message <- m

	}
	fmt.Println("Insert message: ")
	if scanner.Err() != nil {
		log.Println("Error from stdin")
	}

	if localErr != nil {
		log.Println("Error in sending to node")
	}

}
