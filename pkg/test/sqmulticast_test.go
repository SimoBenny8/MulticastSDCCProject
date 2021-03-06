package test

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/SQMulticast"
	client2 "github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/testUtil"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestOneToManySQ(t *testing.T) {
	var localErr error
	var connections []*client2.Client
	var wg sync.WaitGroup

	delay := 1000
	numNode := 3
	messages := [][]byte{[]byte("message"), []byte("in"), []byte("order")}
	connections = make([]*client2.Client, 3)
	for i := 0; i < numNode; i++ {
		connections[i] = testUtil.FakeConnect("Node" + strconv.Itoa(i))
	}

	n := rand.Intn(len(connections))

	node := new(SQMulticast.NodeForSq)
	node.NodeId = uint(rand.Intn(100))
	node.Connections = connections
	node.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
	node.MyConn = connections[0]
	node.LocalTimestamp = 0
	node.ProcessingQueue = make([]*SQMulticast.MessageSeq, 0, 100)

	SQMulticast.AppendNodes(*node)
	go SQMulticast.DeliverMsg(delay, node.NodeId)

	node2 := new(SQMulticast.NodeForSq)
	node2.NodeId = uint(rand.Intn(100))
	node2.Connections = connections
	node2.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
	node2.MyConn = connections[1]
	node2.LocalTimestamp = 0
	node2.ProcessingQueue = make([]*SQMulticast.MessageSeq, 0, 100)

	SQMulticast.AppendNodes(*node2)
	go SQMulticast.DeliverMsg(delay, node2.NodeId)

	node3 := new(SQMulticast.NodeForSq)
	node3.NodeId = uint(rand.Intn(100))
	node3.Connections = connections
	node3.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
	node3.MyConn = connections[2]
	node3.LocalTimestamp = 0
	node3.ProcessingQueue = make([]*SQMulticast.MessageSeq, 0, 100)

	SQMulticast.AppendNodes(*node3)
	go SQMulticast.DeliverMsg(delay, node3.NodeId)

	seq := new(SQMulticast.Sequencer)
	seq.Node = *node
	seq.SeqPort = connections[n]
	seq.Connections = connections
	seq.LocalTimestamp = 0
	seq.MessageQueue = make([]SQMulticast.MessageSeq, 0, 100)
	SQMulticast.SetSequencer(*seq)

	log.Println("Sequencer is", seq.SeqPort.Connection.Target())

	go SQMulticast.DeliverSeq(delay)

	//caso invio al sequencer da un nodo generico
	for i := range messages {
		wg.Add(1)
		for j := range node.Connections {
			if node.Connections[j].Connection.Target() == seq.SeqPort.Connection.Target() {
				//caso invio al sequencer da un nodo generico
				md := make(map[string]string)
				md[util.TYPEMC] = util.SQMULTICAST
				md[util.TYPENODE] = util.SEQUENCER //a chi arriva
				md[util.MESSAGEID] = SQMulticast.RandSeq(5)
				go func() {
					localErr = node.Connections[j].Send(md, []byte(""), messages[i], nil, delay)
					if localErr != nil {
						t.Errorf("SendPacket failed: %v", localErr)
						return

					}
					wg.Done()
				}()

			}
		}
		wg.Wait()

	}

	//assertion
	time.Sleep(time.Second * 30)
	nodes := SQMulticast.GetDeliverNodes()
	for i := range nodes {
		assert.Equal(t, 3, len(nodes[i].DeliverQueue))
	}
	assert.Equal(t, node.DeliverQueue, node2.DeliverQueue)
	assert.Equal(t, node2.DeliverQueue, node3.DeliverQueue)
	assert.Equal(t, node3.DeliverQueue, node.DeliverQueue)

}

func TestManyToManySQ(t *testing.T) {
	var localErr error
	var connections []*client2.Client
	var wg sync.WaitGroup

	delay := 1000
	numNode := 3
	messages := [][]byte{[]byte("message"), []byte("in"), []byte("order")}
	connections = make([]*client2.Client, 3)
	for i := 0; i < numNode; i++ {
		connections[i] = testUtil.FakeConnect("Node" + strconv.Itoa(i))
	}

	n := rand.Intn(len(connections))

	node := new(SQMulticast.NodeForSq)
	node.NodeId = uint(rand.Intn(100))
	node.Connections = connections
	node.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
	node.MyConn = connections[0]
	node.LocalTimestamp = 0
	node.ProcessingQueue = make([]*SQMulticast.MessageSeq, 0, 100)

	SQMulticast.AppendNodes(*node)
	go SQMulticast.DeliverMsg(delay, node.NodeId)

	node2 := new(SQMulticast.NodeForSq)
	node2.NodeId = uint(rand.Intn(1000))
	node2.Connections = connections
	node2.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
	node2.MyConn = connections[1]
	node2.LocalTimestamp = 0
	node2.ProcessingQueue = make([]*SQMulticast.MessageSeq, 0, 100)

	SQMulticast.AppendNodes(*node2)
	go SQMulticast.DeliverMsg(delay, node2.NodeId)

	node3 := new(SQMulticast.NodeForSq)
	node3.NodeId = uint(rand.Intn(1000))
	node3.Connections = connections
	node3.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
	node3.MyConn = connections[2]
	node3.LocalTimestamp = 0
	node3.ProcessingQueue = make([]*SQMulticast.MessageSeq, 0, 100)

	SQMulticast.AppendNodes(*node3)
	go SQMulticast.DeliverMsg(delay, node3.NodeId)

	seq := new(SQMulticast.Sequencer)
	seq.Node = *node
	seq.SeqPort = connections[n]
	seq.Connections = connections
	seq.LocalTimestamp = 0
	seq.MessageQueue = make([]SQMulticast.MessageSeq, 0, 100)
	SQMulticast.SetSequencer(*seq)
	log.Println("Sequencer is", seq.SeqPort.Connection.Target())

	go SQMulticast.DeliverSeq(delay)
	//caso invio al sequencer da un nodo generico

	for i := range messages {
		wg.Add(3)
		for j := range node.Connections {
			if node.Connections[j].Connection.Target() == seq.SeqPort.Connection.Target() {
				//caso invio al sequencer da un nodo generico
				md := make(map[string]string)
				md[util.TYPEMC] = util.SQMULTICAST
				md[util.TYPENODE] = util.SEQUENCER //a chi arriva
				md[util.MESSAGEID] = SQMulticast.RandSeq(5)
				go func() {
					localErr = node.Connections[j].Send(md, []byte(""), messages[i], nil, delay)
					if localErr != nil {
						t.Errorf("SendPacket failed: %v", localErr)
						return

					}
					wg.Done()
				}()
			}
		}

		for j := range node2.Connections {
			if node2.Connections[j].Connection.Target() == seq.SeqPort.Connection.Target() {
				//caso invio al sequencer da un nodo generico
				md := make(map[string]string)
				md[util.TYPEMC] = util.SQMULTICAST
				md[util.TYPENODE] = util.SEQUENCER //a chi arriva
				md[util.MESSAGEID] = SQMulticast.RandSeq(5)
				go func() {
					localErr = node2.Connections[j].Send(md, []byte(""), messages[i], nil, delay)
					if localErr != nil {
						t.Errorf("SendPacket failed: %v", localErr)
						return

					}
					wg.Done()
				}()
			}
		}

		for j := range node3.Connections {
			if node3.Connections[j].Connection.Target() == seq.SeqPort.Connection.Target() {
				//caso invio al sequencer da un nodo generico
				md := make(map[string]string)
				md[util.TYPEMC] = util.SQMULTICAST
				md[util.TYPENODE] = util.SEQUENCER //a chi arriva
				md[util.MESSAGEID] = SQMulticast.RandSeq(5)
				go func() {
					localErr = node3.Connections[j].Send(md, []byte(""), messages[i], nil, delay)
					if localErr != nil {
						t.Errorf("SendPacket failed: %v", localErr)
						return

					}
					wg.Done()
				}()
			}
		}
		wg.Wait()

	}

	//assertion
	time.Sleep(time.Second * 120)
	nodes := SQMulticast.GetDeliverNodes()
	for i := range nodes {
		assert.Equal(t, 9, len(nodes[i].DeliverQueue))
	}
	assert.Equal(t, node.DeliverQueue, node2.DeliverQueue)
	assert.Equal(t, node2.DeliverQueue, node3.DeliverQueue)
	assert.Equal(t, node3.DeliverQueue, node.DeliverQueue)

}
