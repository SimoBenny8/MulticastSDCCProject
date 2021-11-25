package client

import (
	"context"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
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

//create connection between nodes
func Connect(address string) *Client {
	d := rand.Intn(1000) + 1000
	time.Sleep(time.Duration(d) * time.Millisecond)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(address, opts)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Connected client: ", cc.Target())
	c := new(Client)
	c.Client = rpc.NewPacketServiceClient(cc)
	c.Connection = cc
	return &Client{c.Client, c.Connection}
}

//method to send message
func (c *Client) Send(messageMetadata map[string]string, header []byte, payload []byte, respChannel chan []byte, delayParam int) error {
	var wg sync.Mutex
	log.Println("Sender: ", c.Connection.Target())
	wg.Lock()
	delay(&wg, delayParam) //5000 milliseconds
	md := metadata.New(messageMetadata)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err := c.Client.SendPacket(ctx, &rpc.Packet{Header: header, Message: payload})
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

func delay(wg *sync.Mutex, value int) {
	d := rand.Intn(value) + 1000
	time.Sleep(time.Duration(d) * time.Millisecond)
	defer wg.Unlock()
}
