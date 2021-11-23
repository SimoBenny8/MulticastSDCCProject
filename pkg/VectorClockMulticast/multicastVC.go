package VectorClockMulticast

import (
	"bytes"
	"encoding/json"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
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

func (node *NodeVC) InitLocalTimestamp(wg *sync.Mutex, numNode int) {
	node.Timestamp = make([]int, numNode)
	for i := range node.Timestamp {
		node.Timestamp[i] = 0
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

func SendMessageToAll(m *rpc.Packet, port uint, nodeId uint, delay int) {

	pos := checkPositionNode(nodeId)
	var wg sync.WaitGroup
	Nodes[pos].Timestamp[Nodes[pos].MyNode] += 1

	mex := &MessageVectorTimestamp{Address: port, OPacket: *m, Timestamp: Nodes[pos].Timestamp, Id: RandSeq(5)}
	//log.Println("Receiver:", md[util.RECEIVER])
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
	md[util.NODEID] = strconv.Itoa(int(nodeId))
	for i := range Nodes[pos].Connections {
		wg.Add(1)
		ind := i

		go func() {
			err = Nodes[pos].Connections[ind].Send(md, []byte(Nodes[pos].Connections[ind].Connection.Target()+":"+string(m.Header)), b, nil, delay)
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

func (node *NodeVC) ReceiveMessage(message *rpc.Packet) {
	//var err error
	//pos := checkPositionNode(nodeId)
	var wg sync.Mutex
	wg.Lock()
	mt := DecodeMsg(message, &wg)
	//messageT := *mt

	log.Println("Original Message: ", string(mt.OPacket.Message), "Timestamp: ", mt.Timestamp)
	wg.Lock()
	node.AddToProcessingQueue(mt, &wg)

}

func DecodeMsg(message *rpc.Packet, wg *sync.Mutex) *MessageVectorTimestamp {
	// Decode (receive) the value.
	defer wg.Unlock()
	var m MessageVectorTimestamp
	log.Println("decodifica del messaggio")
	var err error
	message.Message = bytes.TrimPrefix(message.Message, []byte("\xef\xbb\xbf"))
	err = json.Unmarshal(message.Message, &m)
	if err != nil {
		log.Printf("error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
	}
	log.Println("messaggio codificato: ", string(m.OPacket.Message))
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return &m

}

func HasSameNumberMessage(mex *MessageVectorTimestamp, nodeId uint) bool {
	pos := checkPositionNode(nodeId)
	resp := false
	for i := range Nodes[pos].Connections {
		if i != int(Nodes[pos].MyNode) {
			if mex.Timestamp[i] <= Nodes[pos].Timestamp[i] {
				log.Println("mex: ", mex.Timestamp[i], " local: ", Nodes[pos].Timestamp[i])
				resp = true
			}
		}
	}
	return resp
}

func indexSender(mex *MessageVectorTimestamp, nodeId uint) int {
	pos := checkPositionNode(nodeId)
	var index int
	for i := range Nodes[pos].Connections {
		if strings.Contains(Nodes[pos].Connections[i].Connection.Target(), strconv.Itoa(int(mex.Address))) {
			index = i
			break
		}
	}
	return index
}

func nextMessageTimestamp(mex *MessageVectorTimestamp, nodeId uint) bool {
	pos := checkPositionNode(nodeId)
	index := indexSender(mex, nodeId)
	if mex.Timestamp[index] == (Nodes[pos].Timestamp[index] + 1) {
		log.Println("next message timestamp 1: ", mex.Timestamp[index], " 2: ", Nodes[pos].Timestamp[index]+1)
		return true
	} else if index == int(Nodes[pos].MyNode) {
		//caso deliver stesso nodo
		return true
	}
	return false
}

func Deliver(nodeId uint, delay int) {
	rand.Seed(time.Now().UnixNano())
	pos := checkPositionNode(nodeId)
	for {
		if len(Nodes[pos].ProcessingMessage) > 0 && nextMessageTimestamp(&Nodes[pos].ProcessingMessage[0], Nodes[pos].NodeId) && HasSameNumberMessage(&Nodes[pos].ProcessingMessage[0], Nodes[pos].NodeId) {
			message := Nodes[pos].Dequeue()
			var wg sync.Mutex
			wg.Lock()
			var LocalErr error
			md := make(map[string]string)
			md[util.TYPEMC] = util.VCMULTICAST
			md[util.ACK] = util.FALSE
			md[util.TIMESTAMPMESSAGE] = util.EMPTY
			md[util.DELIVER] = util.TRUE
			md[util.NODEID] = strconv.Itoa(int(Nodes[pos].NodeId))
			md[util.RECEIVER] = Nodes[pos].MyConn.Connection.Target()
			index := indexSender(&message, Nodes[pos].NodeId)
			for i := range message.Timestamp {
				Nodes[pos].Timestamp[i] = int(math.Max(float64(message.Timestamp[i]), float64(Nodes[pos].Timestamp[i])))
				//log.Println("timestamp locale: ", localTimestamp)
			}
			Nodes[pos].Timestamp[index] += 1
			b, err := json.Marshal(&message)
			if err != nil {
				log.Printf("Error marshalling: %s", err)
				return
			}
			go func(wg *sync.Mutex) {
				defer wg.Unlock()
				LocalErr = Nodes[pos].MyConn.Send(md, []byte(Nodes[pos].MyConn.Connection.Target()+":"), b, nil, delay)
				if LocalErr != nil {
					log.Println(LocalErr.Error())
					Nodes[pos].Timestamp[index] -= 1
					var wg2 sync.Mutex
					wg2.Lock()
					Nodes[pos].AddToProcessingQueue(&message, &wg2)
				}
			}(&wg)

		}
	}

}
