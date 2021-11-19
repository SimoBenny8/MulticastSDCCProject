package test

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/MulticastScalarClock"
	client2 "github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/testUtil"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestOneToManySC(t *testing.T) {
	//rand.Seed(time.Now().UnixNano())

	var connections []*client2.Client
	var wg sync.WaitGroup
	//var port uint

	delay := 5
	//port = 1
	numNode := 3
	messages := [][]byte{[]byte("message"), []byte("in"), []byte("order")}
	connections = make([]*client2.Client, 3)
	for i := 0; i < numNode; i++ {
		connections[i] = testUtil.FakeConnect("Node" + strconv.Itoa(i))
	}

	node := new(MulticastScalarClock.NodeSC)
	node.NodeId = uint(rand.Intn(100))
	node.Connections = connections
	node.ProcessingMessages = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node.MyConn = connections[0]
	node.ReceivedMessage = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node.Timestamp = 0
	node.OrderedAck = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node.OtherTs = make([]MulticastScalarClock.OtherTimestamp, 0, 100)
	node.DeliverQueue = make(MulticastScalarClock.OrderedMessages, 0, 100)

	MulticastScalarClock.AppendNodes(*node)

	go MulticastScalarClock.Receive(uint(1), node.NodeId, delay)
	go MulticastScalarClock.Deliver(len(connections), node.NodeId, delay)

	node2 := new(MulticastScalarClock.NodeSC)
	node2.NodeId = uint(rand.Intn(200))
	node2.Connections = connections
	node2.ProcessingMessages = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node2.MyConn = connections[1]
	node2.ReceivedMessage = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node2.Timestamp = 0
	node2.OrderedAck = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node2.OtherTs = make([]MulticastScalarClock.OtherTimestamp, 0, 100)
	node2.DeliverQueue = make(MulticastScalarClock.OrderedMessages, 0, 100)

	MulticastScalarClock.AppendNodes(*node2)

	go MulticastScalarClock.Receive(uint(2), node2.NodeId, delay)
	go MulticastScalarClock.Deliver(len(connections), node2.NodeId, delay)

	node3 := new(MulticastScalarClock.NodeSC)
	node3.NodeId = uint(rand.Intn(200))
	node3.Connections = connections
	node3.ProcessingMessages = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node3.MyConn = connections[2]
	node3.ReceivedMessage = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node3.Timestamp = 0
	node3.OrderedAck = make(MulticastScalarClock.OrderedMessages, 0, 100)
	node3.OtherTs = make([]MulticastScalarClock.OtherTimestamp, 0, 100)
	node3.DeliverQueue = make(MulticastScalarClock.OrderedMessages, 0, 100)

	MulticastScalarClock.AppendNodes(*node3)

	go MulticastScalarClock.Receive(uint(3), node3.NodeId, delay)
	go MulticastScalarClock.Deliver(len(connections), node3.NodeId, delay)

	for i := range messages {
		wg.Add(1)
		go func() {
			message := &MulticastScalarClock.MessageTimestamp{Address: uint(1), OPacket: rpc.Packet{Message: messages[i]}, Timestamp: MulticastScalarClock.GetTimestamp(node.NodeId), Id: MulticastScalarClock.RandSeq(5)}
			MulticastScalarClock.SendMessageToAll(message, node.NodeId, delay)
			wg.Done()
		}()

		wg.Wait()
	}

	time.Sleep(time.Second * 80)
	nodes := MulticastScalarClock.GetNodes()
	for i := range nodes {
		assert.Equal(t, 3, len(nodes[i].DeliverQueue))
	}
	assert.Equal(t, nodes[0].DeliverQueue, nodes[1].DeliverQueue)
	assert.Equal(t, nodes[1].DeliverQueue, nodes[2].DeliverQueue)
	assert.Equal(t, nodes[2].DeliverQueue, nodes[0].DeliverQueue)

}
