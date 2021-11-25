package VectorClockMulticast

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"log"
	"sync"
)

type VectorMessages []MessageVectorTimestamp

type NodeVC struct {
	NodeId            uint
	Connections       []*client.Client
	DeliverQueue      VectorMessages
	MyConn            *client.Client
	Timestamp         []int
	MyNode            int32
	ProcessingMessage VectorMessages
}

func removeForProcessingMessages(slice VectorMessages, s int) VectorMessages {
	return append(slice[:s], slice[s+1:]...)
}

func (node *NodeVC) AddToProcessingQueue(m *MessageVectorTimestamp, wg *sync.Mutex) {
	defer wg.Unlock()
	log.Println("Message added in queue")
	node.ProcessingMessage = append(node.ProcessingMessage, *m)
}

func (node *NodeVC) Dequeue() MessageVectorTimestamp {
	log.Println("Get the first message of the queue")
	m := node.ProcessingMessage[0]
	node.ProcessingMessage = removeForProcessingMessages(node.ProcessingMessage, 0)
	return m
}
