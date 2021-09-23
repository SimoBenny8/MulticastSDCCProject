package main

import (
	"MulticastSDCCProject/pkg"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	//var port uint
	port := flag.Uint("port", 8090, "server port number")

	var err error
	go func() {

		err = pkg.RunServer(*port)
	}()
	if err != nil {
		log.Println("Error in connection node")
		return
	}

	var localErr error
	var localErrCh2 error
	var localErrCh3 error

	c1 := pkg.Connect("localhost:8090", 1)
	c2 := pkg.Connect("localhost:8091", 1)
	c3 := pkg.Connect("localhost:8092", 1)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			localErr = c1.Send(make(map[string]string), scanner.Bytes())
			localErrCh2 = c2.Send(make(map[string]string), scanner.Bytes())
			localErrCh3 = c3.Send(make(map[string]string), scanner.Bytes())
		}

		if scanner.Err() != nil {
			// Handle error.
		}

		if localErr != nil {
			log.Println("Error in sending to same node")
		}

		if localErrCh2 != nil {
			log.Println("Error in sending to node 2")
		}

		if localErrCh3 != nil {
			log.Println("Error in sending to node 3")
		}
	}

}
