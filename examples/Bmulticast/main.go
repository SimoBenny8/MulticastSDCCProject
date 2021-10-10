package main

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
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

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Insert message: ")
		for scanner.Scan() {

			for i := range connections {
				md := make(map[string]string)
				md["TypeMulticast"] = "BMulticast"
				localErr = connections[i].Send(md, scanner.Bytes(), respChannel)
				result := <-respChannel
				fmt.Println(string(result)) //problema ack implosion
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
