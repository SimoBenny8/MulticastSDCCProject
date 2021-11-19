package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/server"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/pool"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"log"
	"os"
	"strings"
)

func main() {

	//var port uint
	delay := flag.Uint("delay", 1000, "delay for sending operations (ms)")
	port := flag.Uint("port", 8090, "server port number")
	group := flag.String("groupPort", "8090,8091,8092", "defining group port")

	flag.Parse()

	var err error
	go func() {

		err = server.RunServer(*port, server.Register)
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
	pool.Pool.InitThreadPool(connections, 5, util.BMULTICAST, respChannel, *port, 0, int(*delay))

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
