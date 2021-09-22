package pkg

import (
	"MulticastSDCCProject/pkg/rpc"
	"context"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type Client struct {
	Client     rpc.PacketServiceClient
	Connection *grpc.ClientConn
}

//connect with delay
func Connect(address string, delay uint) *Client {
	if delay > 0 {
		time.Sleep(time.Second * time.Duration(delay))
	}
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
func (c *Client) Send(messageHeader map[string]string, payload []byte) error {
	log.Println("Sending message from ", c.Connection.Target())
	_, err := c.Client.SendPacket(context.Background(), &rpc.Packet{Message: payload})
	if err != nil {
		log.Println(err.Error())
		return err
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
func ConnectWithWaitGroup(address string, delay uint, wg *sync.WaitGroup) *Client {
	if delay > 0 {
		time.Sleep(time.Second * time.Duration(delay))
	}
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
