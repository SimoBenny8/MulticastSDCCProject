package VectorClockMulticast

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
	"bytes"
	"context"
	"encoding/json"
	"google.golang.org/grpc/metadata"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MessageVectorTimestamp struct {
	Address   uint
	OPacket   rpc.Packet
	Timestamp []int
	Id        string
	Ack       bool
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//var processingMessage []*rpc.Packet
var respChannel chan []byte
var localTimestamp []int
var network bytes.Buffer
var connections []*client.Client
var IndexNode int32
var myNode int32

func InitLocalTimestamp(wg *sync.Mutex, numNode int) {
	localTimestamp = make([]int, numNode)
	for i := range localTimestamp {
		localTimestamp[i] = 0
	}
	wg.Unlock()
}

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetVectorTimestamp() []int {
	return localTimestamp
}

func SetMyNode(node int32) {
	myNode = node
}

func SetConnections(conn []*client.Client) {
	connections = conn
}

func SendMessageToAll(m *rpc.Packet, port uint) {

	var wg sync.WaitGroup
	localTimestamp[myNode] += 1
	mex := &MessageVectorTimestamp{Address: port, OPacket: *m, Timestamp: localTimestamp, Id: RandSeq(5)}

	b, err := json.Marshal(&mex)
	if err != nil {
		log.Printf("Error marshalling: %s", err)
		return
	}
	md := make(map[string]string)
	md[util.TYPEMC] = util.VCMULTICAST
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

func ReceiveMessage(message *rpc.Packet) {
	//var err error
	var wg sync.Mutex
	wg.Lock()
	mt := DecodeMsg(message, &wg)
	//messageT := *mt

	log.Println("Original Message: ", string(mt.OPacket.Message), "Timestamp: ", mt.Timestamp)
	wg.Lock()
	AddToQueue(mt, &wg)

}

func DecodeMsg(message *rpc.Packet, wg *sync.Mutex) *MessageVectorTimestamp {
	// Decode (receive) the value.
	defer wg.Unlock()
	var m MessageVectorTimestamp
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

func HasSameNumberMessage(mex *MessageVectorTimestamp) bool {
	resp := false
	for i := range connections {
		if i != int(myNode) {
			if mex.Timestamp[i] <= localTimestamp[i] {
				log.Println("mex: ", mex.Timestamp[i], " local: ", localTimestamp[i])
				resp = true
			}
		}
	}
	return resp
}

func indexSender(mex *MessageVectorTimestamp) int {
	var index int
	for i := range connections {
		if strings.Contains(connections[i].Connection.Target(), strconv.Itoa(int(mex.Address))) {
			index = i
			break
		}
	}
	return index
}

func nextMessageTimestamp(mex *MessageVectorTimestamp) bool {
	index := indexSender(mex)
	if mex.Timestamp[index] == (localTimestamp[index] + 1) {
		log.Println("next message timestamp 1: ", mex.Timestamp[index], " 2: ", localTimestamp[index]+1)
		return true
	} else if index == int(myNode) {
		//caso deliver stesso nodo
		return true
	}
	return false
}

func Deliver(myConn *client.Client) {
	rand.Seed(time.Now().UnixNano())

	for {
		if len(GetQueue()) > 0 && nextMessageTimestamp(&queue[0]) && HasSameNumberMessage(&queue[0]) {
			message := Dequeue()
			var wg sync.Mutex
			wg.Lock()
			var LocalErr error
			md := make(map[string]string)
			md[util.TYPEMC] = util.VCMULTICAST
			md[util.ACK] = util.FALSE
			md[util.TIMESTAMPMESSAGE] = util.EMPTY
			md[util.DELIVER] = util.TRUE
			index := indexSender(&message)
			for i := range message.Timestamp {
				localTimestamp[i] = int(math.Max(float64(message.Timestamp[i]), float64(localTimestamp[i])))
				//log.Println("timestamp locale: ", localTimestamp)
			}
			localTimestamp[index] += 1
			metaData := metadata.New(md)
			ctx := metadata.NewOutgoingContext(context.Background(), metaData)
			go func(wg *sync.Mutex) {
				defer wg.Unlock()
				_, LocalErr = myConn.Client.SendPacket(ctx, &message.OPacket)
				if LocalErr != nil {
					log.Println(LocalErr.Error())
					localTimestamp[index] -= 1
					var wg2 sync.Mutex
					wg2.Lock()
					AddToQueue(&message, &wg2)
				}
			}(&wg)

		}
	}

}
