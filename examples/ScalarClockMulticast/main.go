package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/MulticastScalarClock"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/server"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/pool"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//var port uint
	delay := flag.Uint("delay", 1000, "delay for sending operations (ms)")
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

	node := new(MulticastScalarClock.NodeSC)
	node.NodeId = uint(rand.Intn(5))
	node.Connections = connections
	node.ProcessingMessages = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node.MyConn = myConn
	node.ReceivedMessage = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node.Timestamp = 0
	node.OrderedAck = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node.OtherTs = make([]MulticastScalarClock.OtherTimestamp, 0, 100)
	node.DeliverQueue = make(MulticastScalarClock.OrderedMessages, 0, 100)

	MulticastScalarClock.AppendNodes(*node)

	pool.Pool.InitThreadPool(connections, 5, util.SCMULTICAST, nil, *port, node.NodeId, int(*delay))
	go MulticastScalarClock.Receive(*port, node.NodeId, int(*delay))
	go MulticastScalarClock.Deliver(len(connections), node.NodeId, int(*delay))

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
