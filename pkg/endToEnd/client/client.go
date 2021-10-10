package client

import (
	"MulticastSDCCProject/pkg/rpc"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Client struct {
	Client     rpc.PacketServiceClient
	Connection *grpc.ClientConn
}

//connect with delay
func Connect(address string) *Client {

	go time.Sleep(time.Duration(rand.Intn(1700) + 5))

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(address, opts)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("connection state ====> ", cc.GetState(), "connected client ", cc.Target())
	c := new(Client)
	c.Client = rpc.NewPacketServiceClient(cc)
	c.Connection = cc
	return &Client{c.Client, c.Connection}
}

//method to send message
func (c *Client) Send(messageMetadata map[string]string, payload []byte, respChannel chan []byte) error {
	log.Println("Sender: ", c.Connection.Target())
	go time.Sleep(time.Duration(rand.Intn(3700) + 5))
	md := metadata.New(messageMetadata)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err := c.Client.SendPacket(ctx, &rpc.Packet{Message: payload})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	//invio ack
	if respChannel != nil {
		ack := "Message received: " + string(payload)
		respChannel <- []byte(ack)
	}

	return err
}

//method to close connection
func (c *Client) CloseConnection() error {
	err := c.Connection.Close()
	if err != nil {
		return err
	}
	log.Println("Connection closed")
	return nil
}

//connect with delay
func ConnectWithWaitGroup(address string, wg *sync.WaitGroup) *Client {

	go time.Sleep(time.Duration(rand.Intn(700) + 5))

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(address, opts)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("connection state ====> ", cc.GetState(), "connected client ", cc.Target())
	c := new(Client)
	c.Client = rpc.NewPacketServiceClient(cc)
	c.Connection = cc
	wg.Done()
	return &Client{c.Client, c.Connection}
}
