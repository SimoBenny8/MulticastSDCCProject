package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/VectorClockMulticast"
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
	"sync"
	"time"
)

func main() {
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
	var wg sync.Mutex
	var myNode int32

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

	for i := range connections {
		if strings.Contains(connections[i].Connection.Target(), strconv.Itoa(int(*port))) {
			myNode = int32(i)
			log.Println("my index: ", myNode)
		}
	}

	node := new(VectorClockMulticast.NodeVC)
	wg.Lock()
	node.InitLocalTimestamp(&wg, len(connections))
	node.NodeId = uint(rand.Intn(5000))
	node.Connections = connections
	node.DeliverQueue = make(VectorClockMulticast.VectorMessages, 0, 100)
	node.MyConn = myConn
	node.MyNode = myNode
	node.ProcessingMessage = make(VectorClockMulticast.VectorMessages, 0, 100)

	//wg.Lock()
	VectorClockMulticast.AppendNodes(*node)

	pool.Pool.InitThreadPool(connections, 5, util.VCMULTICAST, nil, *port, node.NodeId, int(*delay))
	go VectorClockMulticast.Deliver(node.NodeId, int(*delay))
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
