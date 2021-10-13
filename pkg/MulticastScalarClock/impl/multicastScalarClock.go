package impl

import (
	"MulticastSDCCProject/pkg/SQMulticast"
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

type MessageTimestamp struct {
	Address   uint
	OPacket   rpc.Packet
	Timestamp int32
	Id        string
	Ack       bool
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ProcessingMessage []MessageTimestamp
var respChannel chan []byte
var timestamp int32
var network bytes.Buffer
var connections []*client.Client

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetTimestamp() int32 {
	return timestamp
}

func SendMessageToAll(message *MessageTimestamp, c []*client.Client) {

	message.Timestamp += 1
	timestamp += 1
	connections = c

	b, err := json.Marshal(&message)
	if err != nil {
		fmt.Printf("Error marshalling: %s", err)
		return
	}

	//respChannel = make(chan []byte,1)
	md := make(map[string]string)
	md[SQMulticast.TYPEMC] = SQMulticast.SCMULTICAST
	md[SQMulticast.ACK] = SQMulticast.FALSE
	md[SQMulticast.TIMESTAMPMESSAGE] = SQMulticast.EMPTY
	for i := range connections {
		err = connections[i].Send(md, b, respChannel)
		result := <-respChannel
		log.Println("ack: ", string(result))

		if err != nil {
			log.Fatal("error during send message")
		}
	}

}

func ReceiveMessage(message *rpc.Packet) {
	var err error
	mt := DecodeMsg(message)
	timestamp = int32(math.Max(float64(mt.Timestamp), float64(timestamp)))
	timestamp += 1
	log.Println("Original Message: ", string(mt.OPacket.Message), "Timestamp: ", mt.Timestamp)

	AddToQueue(mt)
	mt.Timestamp = timestamp
	b, err := json.Marshal(&mt)
	if err != nil {
		fmt.Printf("Error marshalling: %s", err)
		return
	}
	//respChannel = make(chan []byte,1)
	md := make(map[string]string)
	md[SQMulticast.TYPEMC] = SQMulticast.SCMULTICAST
	md[SQMulticast.ACK] = SQMulticast.TRUE
	md[SQMulticast.TIMESTAMPMESSAGE] = string(GetQueue()[0].Timestamp) //timestamp del primo messaggio in coda
	for i := range connections {
		err = connections[i].Send(md, b, nil)
		if err != nil {
			log.Fatal("error during ack message")
		}
	}

}

func DecodeMsg(message *rpc.Packet) *MessageTimestamp {
	// Decode (receive) the value.
	var m MessageTimestamp
	log.Println("decodifica del messaggio")
	var err error
	if err = json.Unmarshal(message.Message, &m); err != nil {
		panic(err)
	}
	log.Println("messaggio codificato: ", string(m.OPacket.Message))
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return &m

}

func Deliver() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	for {
		if len(GetQueue()) > 0 {
			message := Dequeue()
			for i := range connections {
				wg.Add(1)
				index := i
				go func() {
					defer wg.Done()
					md := make(map[string]string)
					md[SQMulticast.TYPEMC] = SQMulticast.SCMULTICAST
					md[SQMulticast.ACK] = SQMulticast.FALSE
					md[SQMulticast.TIMESTAMPMESSAGE] = string(GetQueue()[0].Timestamp)
					delay := rand.Intn(10700) + 1000
					//log.Println("Delay: ",delay," milliseconds")
					time.Sleep(time.Duration(delay))
					metaData := metadata.New(md)
					ctx := metadata.NewOutgoingContext(context.Background(), metaData)
					var LocalErr error
					_, LocalErr = connections[index].Client.SendPacket(ctx, &message.OPacket)
					if LocalErr != nil {
						log.Println(LocalErr.Error())
					}
				}()
			}
			wg.Wait()
		}
	}
}

/* Workflow:
 1) Invio messaggio in multicast da parte di Pi
2) Pj, che riceve il messaggio, lo mette in coda e riordina la coda in base al timestamp (incrementare prima di inviare)
3) Pj invia un messaggio di ack in multicast
4) Pj effettua la deliver se msg è in testa e nessun processo ha messaggio con timestamp minore o uguale in coda */
