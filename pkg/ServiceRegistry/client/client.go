package client

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/proto"
	"google.golang.org/grpc"
	"log"
)

//function used to connect a node to the server registry
func Connect(address string) (proto.RegistryClient, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("Error in connecting to registry :", err.Error())
		return nil, err
	} else {
		log.Println("Correctly connected to registry at address: ", address)
	}
	return proto.NewRegistryClient(conn), nil
}
