package main

import (
	"MulticastSDCCProject/pkg/SQMulticast"
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main() {
	port := flag.Uint("port", 8090, "server port number")
	group := flag.String("groupPort", "8091,8092,8093", "defining group port")

	flag.Parse()

	var err error

	go func() {
		err = server.RunServer(*port)
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
	}
	if SQMulticast.Connections == nil {
		SQMulticast.Connections = connections
	}

	if SQMulticast.SeqPort == nil {
		//scelgo quale dei nodi è il sequencer(randomicamente)
		n := rand.Intn(len(connections)) //problema: ogni nodo sceglie il proprio sequencer
		SQMulticast.SeqPort = connections[n]
		log.Println("Sequencer is", SQMulticast.SeqPort.Connection.Target())
	}

	//implemento invio messaggio
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Insert message: ")
		for scanner.Scan() {

			for i := range connections {
				if connections[i].Connection.Target() == SQMulticast.SeqPort.Connection.Target() {
					//caso invio al sequencer da un nodo generico
					md := make(map[string]string)
					md[SQMulticast.TYPEMC] = SQMulticast.SQMULTICAST
					md[SQMulticast.TYPENODE] = SQMulticast.MEMBER //da chi arriva
					md[SQMulticast.MESSAGEID] = SQMulticast.RandSeq(5)
					localErr = connections[i].Send(md, scanner.Bytes(), nil)
				} else {
					log.Println("Impossible to send message from Sequencer")
					continue
				}

			}

		}

		if scanner.Err() != nil {
			log.Println("Error from stdin")
		}

		if localErr != nil {
			log.Println("Error in sending to node")
		}

	}

}
