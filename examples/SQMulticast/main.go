package main

import (
	"MulticastSDCCProject/pkg/SQMulticast"
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"MulticastSDCCProject/pkg/pool"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main() {
	port := flag.Uint("port", 8090, "server port number")
	group := flag.String("groupPort", "8090,8091,8092", "defining group port")

	flag.Parse()

	var err error

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
	//myAddress := "localhost:"+ strconv.Itoa(int(*port))

	connections = make([]*client.Client, len(groupArray))
	rand.Seed(int64(len(groupArray)))

	//connessione ai nodi
	for i := range groupArray {
		connections[i] = client.Connect("localhost:" + groupArray[i])
	}

	if SQMulticast.SeqPort == nil {
		//scelgo quale dei nodi è il sequencer(randomicamente)
		n := rand.Intn(len(connections))
		SQMulticast.SeqPort = connections[n]
		log.Println("Sequencer is", SQMulticast.SeqPort.Connection.Target())
	}
	pool.Pool.InitThreadPool(connections, 5, util.SQMULTICAST, nil, *port)
	go SQMulticast.DeliverSeq(connections)

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
