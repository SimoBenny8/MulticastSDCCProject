package VectorClockMulticast

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"log"
	"sync"
)

type VectorMessages []MessageVectorTimestamp

type NodeVC struct {
	NodeId             uint
	Connections        []*client.Client
	DeliverQueue       VectorMessages
	MyConn             *client.Client
	Timestamp          []int
	ProcessingMessages VectorMessages
	myNode             int32
}

// Create
var queue VectorMessages

func init() {
	//	queue = make(VectorMessages, 0, 100)
}

func (node *NodeVC) AddToQueue(m *MessageVectorTimestamp, wg *sync.Mutex) {
	defer wg.Unlock()
	log.Println("messaggio aggiunto in coda")
	node.ProcessingMessages = append(node.ProcessingMessages, *m)
}

func (node *NodeVC) Dequeue() MessageVectorTimestamp {
	log.Println("prendo il primo elemento della coda")
	m := node.ProcessingMessages[0]
	if len(node.ProcessingMessages) > 1 {
		node.ProcessingMessages = node.ProcessingMessages[1:]
	} else {
		node.ProcessingMessages = node.ProcessingMessages[:0]
	}
	return m
}
