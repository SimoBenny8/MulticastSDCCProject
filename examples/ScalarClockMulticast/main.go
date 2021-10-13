package main

import (
	"MulticastSDCCProject/pkg/MulticastScalarClock/impl"
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"MulticastSDCCProject/pkg/rpc"
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
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

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Insert message: ")
		for scanner.Scan() {

			m := &rpc.Packet{Message: scanner.Bytes()}
			mex := &impl.MessageTimestamp{Address: *port, OPacket: *m, Timestamp: impl.GetTimestamp(), Id: impl.RandSeq(5)}
			impl.SendMessageToAll(mex, connections)

		}
		fmt.Println("Insert message: ")
		if scanner.Err() != nil {
			log.Println("Error from stdin")
		}

		if localErr != nil {
			log.Println("Error in sending to node")
		}
	}

}
