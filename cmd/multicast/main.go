package main

import (
	"flag"
	"fmt"
	server2 "github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/server"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/server"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/restApi"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {

	delay := flag.Uint("delay", uint(util.GetEnvIntWithDefault("DELAY", 0)), "delay for sending operations (ms)")
	grpcPort := flag.Uint("rpcPort", uint(util.GetEnvIntWithDefault("GRPC_PORT", 90)), "port number of the grpc server")
	restPort := flag.Uint("restPort", uint(util.GetEnvIntWithDefault("REST_PORT", 80)), "port number of the rest server")
	restPath := flag.String("restPath", util.GetEnvStringWithDefault("REST_PATH", "/multicast/v1"), "path of the rest api")
	numThreads := flag.Uint("nThreads", uint(util.GetEnvIntWithDefault("NUM_THREADS", 1)), "number of threads used to multicast messages")
	verb := flag.Bool("verbose", util.GetEnvBoolWithDefault("VERBOSE", true), "Turn verbose mode on or off.")
	registryAddr := flag.String("registryAddr", "registry:90", "service registry address")
	r := flag.Bool("registry", util.GetEnvBoolWithDefault("REGISTRY", false), "start multicast registry")
	application := flag.Bool("application", util.GetEnvBoolWithDefault("APP", false), "start multicast application")

	var myPort []string
	flag.Parse()
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)

	if *application {
		log.Println("Adding basic communication service to gRPC server")
		services = append(services, server.Register)
	}
	if *r {
		log.Println("Adding multicast registry service to gRPC server")
		services = append(services, server2.Registration)
	}
	log.Println("start services")
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		myPort = strings.SplitAfter(lis.Addr().String(), ".")
		if err != nil {
			return
		}
		log.Printf("Grpc-Server started at %v", lis.Addr().String())

		s := grpc.NewServer()
		for _, grpcService := range services {
			err = grpcService(s)
			if err != nil {
				return
			}

		}

		if err = s.Serve(lis); err != nil {
			return
		}
		wg.Done()
	}()

	if *application {

		wg.Add(1)
		go func() {
			err := restApi.Run(*grpcPort, *restPort, *registryAddr, *restPath, int(*numThreads), *delay, *verb, myPort[len(myPort)-1])
			if err != nil {
				log.Println("Error in running application", err.Error())
				return
			}
			wg.Done()
		}()
	}

	wgChan := make(chan bool)

	go func() {
		wg.Wait()
		wgChan <- true
	}()

	log.Println("App started")
	select {
	case <-wgChan:
	}

}
