package impl

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
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
	Address        uint
	OPacket        rpc.Packet
	Timestamp      int
	Id             string
	Ack            bool
	FirstTsInQueue int
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var processingMessage []*rpc.Packet
var respChannel chan []byte
var timestamp int
var network bytes.Buffer
var connections []*client.Client

func AddingRecevingMex(mex *rpc.Packet) {
	processingMessage = append(processingMessage, mex)
}

func init() {
	processingMessage = make([]*rpc.Packet, 0, 100)
}

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetTimestamp() int {
	return timestamp
}

func SetConnections(c []*client.Client) {
	connections = c
}

func SendMessageToAll(message *MessageTimestamp) {

	var wg sync.WaitGroup
	message.Timestamp += 1
	timestamp += 1

	b, err := json.Marshal(&message)
	if err != nil {
		fmt.Printf("Error marshalling: %s", err)
		return
	}

	//respChannel = make(chan []byte,1)
	md := make(map[string]string)
	md[util.TYPEMC] = util.SCMULTICAST
	md[util.ACK] = util.FALSE
	md[util.TIMESTAMPMESSAGE] = util.EMPTY
	md[util.DELIVER] = util.FALSE
	for i := range connections {
		wg.Add(1)
		ind := i
		go func() {
			err = connections[ind].Send(md, b, nil)
			//result := <-respChannel
			//log.Println("ack: ", string(result))
			if err != nil {
				log.Fatal("error during send message")
			}
			wg.Done()
		}()

	}
	wg.Wait()
}

func DecodeMsg(message *rpc.Packet) *MessageTimestamp {
	// Decode (receive) the value.
	var m MessageTimestamp
	log.Println("decodifica del messaggio")
	var err error

	message.Message = bytes.TrimPrefix(message.Message, []byte("\xef\xbb\xbf"))
	if err = json.Unmarshal(message.Message, &m); err != nil {
		panic(err)
	}
	log.Println("messaggio codificato: ", string(m.OPacket.Message))
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return &m

}

func HasMinimumTimestamp(message *MessageTimestamp) bool {
	var count int

	for i := range otherTs {
		if otherTs[i].otherNodeTimestamp >= message.Timestamp {
			count += 1
		}
	}
	log.Println("count:", count)
	if count == len(connections) {
		EmptyOtherTimestamp(message.Id)
		return true
	}
	return false
}

func Deliver(myConn *client.Client, numConn int) {
	rand.Seed(time.Now().UnixNano())

	for {
		if len(GetQueue()) > 0 && IsCorrectNumberAck(&queue[0], numConn, nil) && HasMinimumTimestamp(&queue[0]) {
			var wg sync.Mutex
			wg.Lock()
			message := Dequeue()
			var LocalErr error
			md := make(map[string]string)
			md[util.TYPEMC] = util.SCMULTICAST
			md[util.ACK] = util.FALSE
			md[util.TIMESTAMPMESSAGE] = util.EMPTY
			md[util.DELIVER] = util.TRUE
			metaData := metadata.New(md)
			ctx := metadata.NewOutgoingContext(context.Background(), metaData)
			go func(wg *sync.Mutex) {
				defer wg.Unlock()
				_, LocalErr = myConn.Client.SendPacket(ctx, &message.OPacket)
				if LocalErr != nil {
					log.Println(LocalErr.Error())
				}

			}(&wg)

		}
	}

}

func Receive(conn []*client.Client, addr uint) {
	var tInQueue int
	for {
		if len(processingMessage) > 0 {
			var wg sync.WaitGroup
			mt := DecodeMsg(processingMessage[0])
			timestamp = int(math.Max(float64(mt.Timestamp), float64(timestamp)))
			timestamp += 1
			log.Println("Original Message: ", string(mt.OPacket.Message), "Timestamp: ", mt.Timestamp)

			AddToQueue(mt)
			if len(GetQueue()) > 0 {
				tInQueue = queue[0].Timestamp
			} else {
				tInQueue = -1
			}
			mt.Timestamp = timestamp
			mt.Ack = true
			mt.Address = addr
			mt.FirstTsInQueue = tInQueue
			b, err := json.Marshal(&mt)
			if err != nil {
				fmt.Printf("Error marshalling: %s", err)
				return
			}

			md := make(map[string]string)
			md[util.TYPEMC] = util.SCMULTICAST
			md[util.ACK] = util.TRUE
			md[util.DELIVER] = util.FALSE
			for i := range conn {
				wg.Add(1)
				ind := i
				go func() {
					defer wg.Done()

					err = conn[ind].Send(md, b, nil)
					if err != nil {
						log.Fatal("error during ack message")
					}
				}()
			}
			if len(processingMessage) > 1 {
				processingMessage = processingMessage[1:]
			} else {
				processingMessage = processingMessage[:0]
			}
			wg.Wait()
		} else {
			continue
		}
	}
}

func IsCorrectNumberAck(message *MessageTimestamp, numCon int, wg *sync.Mutex) bool {
	i := 0
	for r := range orderedAck {
		if orderedAck[r].Id == message.Id {
			i += 1
		}
		log.Println("numero ack: ", i)
	}

	if i == numCon {
		log.Println("raggiunto numero di ack corretto")
		EmptyOrderedAck(message.Id)
		//	wg.Unlock()
		return true

	}
	//wg.Unlock()
	return false
}

/* Workflow:
 1) Invio messaggio in multicast da parte di Pi
2) Pj, che riceve il messaggio, lo mette in coda e riordina la coda in base al timestamp (incrementare prima di inviare)
3) Pj invia un messaggio di ack in multicast
4) Pj effettua la deliver se msg è in testa e nessun processo ha messaggio con timestamp minore o uguale in coda */
