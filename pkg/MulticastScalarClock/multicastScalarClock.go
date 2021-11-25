package MulticastScalarClock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"log"
	"math"
	"math/rand"
	"strconv"
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

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func (node *NodeSC) AddingReceivingMex(mex *rpc.Packet) {
	m := DecodeMsg(mex)
	node.ReceivedMessage = append(node.ReceivedMessage, *m)
}

func GetTimestamp(nodeId uint) int {
	pos := checkPositionNode(nodeId)
	return Nodes[pos].Timestamp
}

//function used to get a random sequence for message ID
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//Send message in multicast
func SendMessageToAll(header []byte, message *MessageTimestamp, nodeId uint, delay int) {

	pos := checkPositionNode(nodeId)
	var wg sync.WaitGroup
	message.Timestamp += 1
	Nodes[pos].Timestamp += 1

	b, err := json.Marshal(&message)
	if err != nil {
		fmt.Printf("Error marshalling: %s", err)
		return
	}

	md := make(map[string]string)
	md[util.TYPEMC] = util.SCMULTICAST
	md[util.ACK] = util.FALSE
	md[util.TIMESTAMPMESSAGE] = util.EMPTY
	md[util.DELIVER] = util.FALSE
	md[util.NODEID] = strconv.Itoa(int(nodeId))
	for i := range Nodes[pos].Connections {
		wg.Add(1)
		ind := i
		go func() {
			err = Nodes[pos].Connections[ind].Send(md, []byte(Nodes[pos].Connections[ind].Connection.Target()+":"+string(header)), b, nil, delay)
			if err != nil {
				log.Fatal("error during send message")
			}
			wg.Done()
		}()

	}
	wg.Wait()
}

//function used to unmarshal a received message
func DecodeMsg(message *rpc.Packet) *MessageTimestamp {
	// Decode (receive) the value.
	var m MessageTimestamp
	var err error

	message.Message = bytes.TrimPrefix(message.Message, []byte("\xef\xbb\xbf"))
	if err = json.Unmarshal(message.Message, &m); err != nil {
		panic(err)
	}
	log.Println("decoded messages: ", string(m.OPacket.Message))
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return &m

}

//function that checks that the message that has to be delivered, has minimum timestamp
func HasMinimumTimestamp(message *MessageTimestamp, nodeId uint) bool {
	var count int
	pos := checkPositionNode(nodeId)
	for i := range Nodes[pos].OtherTs {
		if Nodes[pos].OtherTs[i].otherNodeTimestamp >= message.Timestamp {
			count += 1
		}
	}
	//log.Println("count:", count)
	if count >= len(Nodes[pos].Connections) {
		log.Println("count:", count)
		return true
	}
	return false
}

//function that checks if a message has to be delivered
func Deliver(numConn int, nodeId uint, delay int) {
	rand.Seed(time.Now().UnixNano())
	pos := checkPositionNode(nodeId)
	for {
		if len(Nodes[pos].ProcessingMessages) > 0 && IsCorrectNumberAck(&Nodes[pos].ProcessingMessages[0], numConn, nodeId) && HasMinimumTimestamp(&Nodes[pos].ProcessingMessages[0], nodeId) {
			var wg sync.Mutex
			wg.Lock()
			message := Nodes[pos].Dequeue()
			EmptyOrderedAck(message.Id, nodeId)
			EmptyOtherTimestamp(message.Id, nodeId)
			var LocalErr error
			md := make(map[string]string)
			md[util.TYPEMC] = util.SCMULTICAST
			md[util.ACK] = util.FALSE
			md[util.TIMESTAMPMESSAGE] = util.EMPTY
			md[util.DELIVER] = util.TRUE
			md[util.NODEID] = strconv.Itoa(int(nodeId))
			md[util.RECEIVER] = Nodes[pos].MyConn.Connection.Target()
			b, err := json.Marshal(&message)
			if err != nil {
				log.Printf("Error marshalling: %s", err)
				return
			}
			go func(wg *sync.Mutex) {
				defer wg.Unlock()
				LocalErr = Nodes[pos].MyConn.Send(md, []byte(Nodes[pos].MyConn.Connection.Target()), b, nil, delay)
				if LocalErr != nil {
					log.Println(LocalErr.Error())
				}

			}(&wg)

		}
	}

}

//function that manage a received message
func Receive(addr uint, nodeId uint, delay int) {
	var tInQueue int
	pos := checkPositionNode(nodeId)
	for {
		if len(Nodes[pos].ReceivedMessage) > 0 {
			var wg sync.WaitGroup
			var mutex sync.Mutex
			mt := Nodes[pos].ReceivedMessage[0]
			Nodes[pos].Timestamp = int(math.Max(float64(mt.Timestamp), float64(Nodes[pos].Timestamp)))
			Nodes[pos].Timestamp += 1
			log.Println("Original Message: ", string(mt.OPacket.Message), "Timestamp: ", mt.Timestamp)
			mutex.Lock()
			Nodes[pos].AddToQueue(&mt, &mutex)
			if len(Nodes[pos].ProcessingMessages) > 0 {
				tInQueue = Nodes[pos].ProcessingMessages[0].Timestamp
			} else {
				tInQueue = -1
			}
			mt.Timestamp = Nodes[pos].Timestamp
			mt.Ack = true
			mt.Address = addr
			mt.FirstTsInQueue = tInQueue
			b, err := json.Marshal(&mt)
			if err != nil {
				fmt.Printf("Error marshalling: %s", err)
				return
			}
			Nodes[pos].ReceivedMessage = removeForReceivedMessage(Nodes[pos].ReceivedMessage, 0)
			md := make(map[string]string)
			md[util.TYPEMC] = util.SCMULTICAST
			md[util.ACK] = util.TRUE
			md[util.DELIVER] = util.FALSE
			md[util.NODEID] = strconv.Itoa(int(nodeId))
			for i := range Nodes[pos].Connections {
				wg.Add(1)
				ind := i
				go func() {
					defer wg.Done()
					err = Nodes[pos].Connections[ind].Send(md, []byte(Nodes[pos].Connections[ind].Connection.Target()), b, nil, delay)
					if err != nil {
						log.Fatal("error during ack message")
					}
				}()
			}
			wg.Wait()

		}
	}
}

//function that checks if all ack of a message has arrived
func IsCorrectNumberAck(message *MessageTimestamp, numCon int, nodeId uint) bool {
	i := 0
	pos := checkPositionNode(nodeId)
	for r := range Nodes[pos].OrderedAck {
		if Nodes[pos].OrderedAck[r].Id == message.Id {
			i += 1
		}
	}

	if i == numCon {
		return true

	}
	return false
}
