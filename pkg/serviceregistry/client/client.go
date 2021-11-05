package client

import (
	"MulticastSDCCProject/pkg/serviceregistry/ServiceProto"
	"google.golang.org/grpc"
	"log"
)

func Connect(address string) (ServiceProto.RegistryClient, error) {

	log.Println("Connecting to registry: ", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("Error in connecting to registry :", err.Error())
		return nil, err
	} else {
		log.Println("Correctly connected to registry at address: ", address)
	}
	return ServiceProto.NewRegistryClient(conn), nil
}
