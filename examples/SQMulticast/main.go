package main

import (
	"bufio"
	"flag"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/SQMulticast"
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
)

func main() {
	delay := flag.Uint("delay", 1000, "delay for sending operations (ms)")
	port := flag.Uint("port", 8090, "server port number")
	group := flag.String("groupPort", "8090,8091,8092", "defining group port")

	flag.Parse()

	var err error
	var myConn *client.Client

	go func() {
		err = server.RunServer(*port, server.Register)
	}()
	if err != nil {
		log.Println("Error in connection sequencer")
		return
	}

	var localErr error
	var groupArray []string
	var connections []*client.Client

	groupArray = strings.Split(*group, ",")

	connections = make([]*client.Client, len(groupArray))
	rand.Seed(int64(len(groupArray)))

	//connessione ai nodi
	for i := range groupArray {
		connections[i] = client.Connect("localhost:" + groupArray[i])
		portConn, localErr := strconv.Atoi(groupArray[i])
		if localErr != nil {
			log.Println("Error from atoi")
		}
		if portConn == int(*port) {
			myConn = connections[i]
		}
	}

	n := rand.Intn(len(connections))

	//node.Connections = connections

	node := new(SQMulticast.NodeForSq)
	node.NodeId = uint(rand.Intn(5))
	node.Connections = connections
	node.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
	node.MyConn = myConn

	SQMulticast.AppendNodes(*node)

	seq := new(SQMulticast.Sequencer)
	seq.Node = *node
	seq.SeqPort = connections[n]
	seq.Connections = connections
	seq.LocalTimestamp = 0
	seq.MessageQueue = make([]SQMulticast.MessageSeq, 0, 100)
	SQMulticast.SetSequencer(*seq)
	log.Println("Sequencer is", seq.SeqPort.Connection.Target())

	go SQMulticast.DeliverSeq(int(*delay))
	pool.Pool.InitThreadPool(connections, 5, util.SQMULTICAST, nil, *port, node.NodeId, int(*delay))

	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Insert message: ")
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
