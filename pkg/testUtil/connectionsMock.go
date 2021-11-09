package testUtil

import (
	client2 "MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"MulticastSDCCProject/pkg/rpc"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener
)

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	rpc.RegisterPacketServiceServer(s, &server.Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func FakeConnect(address string) *client2.Client {
	ctx := context.Background()

	cc, err := grpc.DialContext(ctx, address, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}
	//log.Println("connection state ====> ", cc.GetState(), "connected client ", cc.Target())
	c := new(client2.Client)
	c.Client = rpc.NewPacketServiceClient(cc)
	c.Connection = cc
	return &client2.Client{Client: c.Client, Connection: c.Connection}
}
