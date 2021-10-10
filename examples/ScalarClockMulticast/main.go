package main

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	impl2 "MulticastSDCCProject/pkg/multicastScalarClock/impl"
	"MulticastSDCCProject/pkg/rpc"
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
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
	var timestamp int32
	var groupArray []string
	var connections []*client.Client
	var wg sync.WaitGroup

	go func() {

		err = server.RunServer(*port)
	}()
	if err != nil {
		log.Println("Error in connection node")
		return
	}

	groupArray = strings.Split(*group, ",")

	connections = make([]*client.Client, len(groupArray))

	for i := range groupArray {
		connections[i] = client.Connect("localhost:" + groupArray[i])
	}

	//var buffNetwork bytes.Buffer
	timestamp = 0
	impl2.InitQueue() //queueing deliver messages

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Insert message: ")
	for scanner.Scan() {

		m := &rpc.Packet{Message: scanner.Bytes()}
		mex := &impl2.MessageTimestamp{Address: *port, OPacket: *m, Timestamp: timestamp, Id: impl2.RandSeq(5)}
		wg.Add(1)
		impl2.SendMessageToAll(mex, connections, &timestamp, &wg)

		wg.Wait()
	}
	fmt.Println("Insert message: ")
	if scanner.Err() != nil {
		log.Println("Error from stdin")
	}

	if localErr != nil {
		log.Println("Error in sending to node")
	}

}
