package MulticastScalarClock

import (
	"log"
	"sort"
	"sync"
)

type OrderedMessages []MessageTimestamp

func (a OrderedMessages) Len() int           { return len(a) }
func (a OrderedMessages) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }
func (a OrderedMessages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//Order message in queue
func (node *NodeSC) OrderingMessage() {
	sort.Sort(node.ProcessingMessages)
	return
}

func (node *NodeSC) AddToQueue(m *MessageTimestamp, wg *sync.Mutex) {

	log.Println("Message added in queue of node", node.NodeId)
	node.ProcessingMessages = append(node.ProcessingMessages, *m)
	node.OrderingMessage()
	wg.Unlock()
}

//Get a message from queue
func (node *NodeSC) Dequeue() MessageTimestamp {
	log.Println("Get the first message of node:", node.NodeId)
	m := node.ProcessingMessages[0]
	node.ProcessingMessages = removeForProcessingMessages(node.ProcessingMessages, 0)
	return m
}
