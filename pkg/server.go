package pkg

import (
	"MulticastSDCCProject/pkg/rpc"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

// RunServer runs gRPC service to publish Message service
func RunServer(port uint) error {

	netListener, err := getNetListener(port)
	if err == nil {
		log.Println("Succed to listen : ", port)
	}
	// register service
	server := grpc.NewServer()
	svr := rpc.NewMessageGrpcImpl()
	rpc.RegisterPacketServiceServer(server, svr)
	// start the server
	err = server.Serve(netListener)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	} else {
		log.Println("server connected")
	}
	return err
}

func RunServerWithWaitGroup(port uint, wg *sync.WaitGroup) error {

	netListener, err := getNetListener(port)
	if err == nil {
		log.Println("Succed to listen : ", port)
	}
	// register service
	server := grpc.NewServer()
	svr := rpc.NewMessageGrpcImpl()
	rpc.RegisterPacketServiceServer(server, svr)
	// start the server
	err = server.Serve(netListener)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	} else {
		log.Println("server connected")
	}
	wg.Done()
	return err
}

func getNetListener(port uint) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	return lis, err
}
