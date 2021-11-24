package test

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/VectorClockMulticast"
	client2 "github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/testUtil"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"strings"
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

	time.Sleep(time.Second * 40)
	nodes := VectorClockMulticast.GetNodes()
	for i := range nodes {
		assert.Equal(t, 3, len(nodes[i].DeliverQueue))
		assert.Equal(t, messages[0], nodes[i].DeliverQueue[0].OPacket.Message)
		assert.Equal(t, messages[1], nodes[i].DeliverQueue[1].OPacket.Message)
		assert.Equal(t, messages[2], nodes[i].DeliverQueue[2].OPacket.Message)
	}

}

func getPositionMessage(slice VectorClockMulticast.VectorMessages, message []byte) int {
	for i := 0; i < len(slice); i++ {
		if strings.Contains(string(slice[i].OPacket.Message), string(message)) {
			return i
		}
	}
	return -1
}

func TestManyToManyVC(t *testing.T) {
	//rand.Seed(time.Now().UnixNano())

	var connections []*client2.Client
	var wg sync.Mutex
	var wg2 sync.WaitGroup
	//var port uint

	delay := 5
	//port = 1
	numNode := 3
	messages := [][]byte{[]byte("M1"), []byte("M2"), []byte("M3")}
	messages2 := [][]byte{[]byte("R:M1"), []byte("R:M2"), []byte("R:M3")}
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

	m1 := func(wg2 *sync.WaitGroup) {
		for i := range messages {
			time.Sleep(time.Second * 2)
			m := &rpc.Packet{Message: messages[i]}
			VectorClockMulticast.SendMessageToAll(m, 0, node.NodeId, delay)
		}
		wg2.Done()
	}

	m2 := func() {
		for {
			nodes := VectorClockMulticast.GetNodes()
			if len(nodes[1].DeliverQueue) > 0 {
				m := &rpc.Packet{Message: messages2[0]}
				VectorClockMulticast.SendMessageToAll(m, 1, node2.NodeId, delay)
				break
			}
		}
	}

	m3 := func() {
		for {
			nodes := VectorClockMulticast.GetNodes()
			if len(nodes[2].DeliverQueue) > 0 {
				time.Sleep(time.Second)
				m := &rpc.Packet{Message: messages2[1]}
				VectorClockMulticast.SendMessageToAll(m, 2, node3.NodeId, delay)
				break
			}
		}
	}

	go m2()
	go m3()

	wg2.Add(1)
	go m1(&wg2)
	wg2.Wait()

	time.Sleep(time.Second * 150)
	nodes := VectorClockMulticast.GetNodes()
	for i := range nodes {
		assert.Equal(t, 5, len(nodes[i].DeliverQueue))
	}
	assert.True(t, getPositionMessage(nodes[0].DeliverQueue, messages[0]) < getPositionMessage(nodes[0].DeliverQueue, messages2[0]))
	assert.True(t, getPositionMessage(nodes[1].DeliverQueue, messages[0]) < getPositionMessage(nodes[1].DeliverQueue, messages2[0]))
	assert.True(t, getPositionMessage(nodes[2].DeliverQueue, messages[0]) < getPositionMessage(nodes[2].DeliverQueue, messages2[0]))
	assert.True(t, getPositionMessage(nodes[0].DeliverQueue, messages2[0]) < getPositionMessage(nodes[0].DeliverQueue, messages2[1]))
	assert.True(t, getPositionMessage(nodes[1].DeliverQueue, messages2[0]) < getPositionMessage(nodes[1].DeliverQueue, messages2[1]))
	assert.True(t, getPositionMessage(nodes[2].DeliverQueue, messages2[0]) < getPositionMessage(nodes[2].DeliverQueue, messages2[1]))

}
