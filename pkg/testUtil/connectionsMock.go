package testUtil

import (
	"context"
	client2 "github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/server"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
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

//Creates a BufCon connection between nodes
func FakeConnect(address string) *client2.Client {
	ctx := context.Background()

	cc, err := grpc.DialContext(ctx, address, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}
	c := new(client2.Client)
	c.Client = rpc.NewPacketServiceClient(cc)
	c.Connection = cc
	return &client2.Client{Client: c.Client, Connection: c.Connection}
}
