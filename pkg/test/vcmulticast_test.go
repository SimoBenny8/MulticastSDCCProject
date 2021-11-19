package test

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/VectorClockMulticast"
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

func TestOneToManyVC(t *testing.T) {
	//rand.Seed(time.Now().UnixNano())

	var connections []*client2.Client
	var wg sync.Mutex
	var wg2 sync.WaitGroup
	//var port uint

	delay := 5
	//port = 1
	numNode := 3
	messages := [][]byte{[]byte("message"), []byte("in"), []byte("order")}
	connections = make([]*client2.Client, 3)
	for i := 0; i < numNode; i++ {
		connections[i] = testUtil.FakeConnect("Node" + strconv.Itoa(i))
	}

	node := new(VectorClockMulticast.NodeVC)
	node.NodeId = uint(rand.Intn(100))
	node.Connections = connections
	node.MyConn = connections[0]
	wg.Lock()
	node.InitLocalTimestamp(&wg, 3)
	node.DeliverQueue = make(VectorClockMulticast.VectorMessages, 0, 100)
	node.ProcessingMessage = make(VectorClockMulticast.VectorMessages, 0, 100)
	node.MyNode = 0

	wg.Lock()
	VectorClockMulticast.AppendNodes(*node, &wg)
	go VectorClockMulticast.Deliver(node.NodeId, delay)

	node2 := new(VectorClockMulticast.NodeVC)
	node2.NodeId = uint(rand.Intn(100))
	node2.Connections = connections
	node2.MyConn = connections[1]
	wg.Lock()
	node2.InitLocalTimestamp(&wg, 3)
	node2.DeliverQueue = make(VectorClockMulticast.VectorMessages, 0, 100)
	node2.ProcessingMessage = make(VectorClockMulticast.VectorMessages, 0, 100)
	node2.MyNode = 1

	wg.Lock()
	VectorClockMulticast.AppendNodes(*node2, &wg)
	go VectorClockMulticast.Deliver(node2.NodeId, delay)

	node3 := new(VectorClockMulticast.NodeVC)
	node3.NodeId = uint(rand.Intn(100))
	node3.Connections = connections
	node3.MyConn = connections[2]
	wg.Lock()
	node3.InitLocalTimestamp(&wg, 3)
	node3.DeliverQueue = make(VectorClockMulticast.VectorMessages, 0, 100)
	node3.ProcessingMessage = make(VectorClockMulticast.VectorMessages, 0, 100)
	node3.MyNode = 2

	wg.Lock()
	VectorClockMulticast.AppendNodes(*node3, &wg)
	go VectorClockMulticast.Deliver(node3.NodeId, delay)

	for i := range messages {
		wg2.Add(1)
		go func() {
			m := &rpc.Packet{Message: messages[i]}
			VectorClockMulticast.SendMessageToAll(m, 0, node.NodeId, delay)
			wg2.Done()
		}()

		wg2.Wait()
	}

	time.Sleep(time.Second * 30)
	nodes := VectorClockMulticast.GetNodes()
	for i := range nodes {
		assert.Equal(t, 3, len(nodes[i].DeliverQueue))
		assert.Equal(t, messages[0], nodes[i].DeliverQueue[0].OPacket.Message)
		assert.Equal(t, messages[1], nodes[i].DeliverQueue[1].OPacket.Message)
		assert.Equal(t, messages[2], nodes[i].DeliverQueue[2].OPacket.Message)

	}

}
