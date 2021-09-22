package main

import (
	"MulticastSDCCProject/pkg"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

	var err error
	go func() {

		err = pkg.RunServer(8092)
	}()

	if err != nil {
		log.Println("Error in connection node 3")
		return
	}

	c1 := pkg.Connect("localhost:8090", 1)
	c2 := pkg.Connect("localhost:8091", 1)

	var localErr error
	var localErrCh2 error

	for {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			localErr = c1.Send(make(map[string]string), scanner.Bytes())
			localErrCh2 = c2.Send(make(map[string]string), scanner.Bytes())
		}

		if scanner.Err() != nil {
			// Handle error.
		}

		if localErr != nil {
			log.Println("Error in sending to node 1")
		}

		if localErrCh2 != nil {
			log.Println("Error in sending to node 3")
		}
	}

}
