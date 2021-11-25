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

//func che ordina i messaggi in coda
func (node *NodeSC) OrderingMessage() {
	sort.Sort(node.ProcessingMessages)
	//log.Println("ordering messages", node.ProcessingMessages[0].Timestamp)
	return
}

func (node *NodeSC) AddToQueue(m *MessageTimestamp, wg *sync.Mutex) {

	log.Println("messaggio aggiunto in coda", node.NodeId)
	node.ProcessingMessages = append(node.ProcessingMessages, *m)
	node.OrderingMessage()
	wg.Unlock()
}

func (node *NodeSC) Dequeue() MessageTimestamp {
	log.Println("prendo il primo elemento della coda", node.NodeId)
	m := node.ProcessingMessages[0]
	node.ProcessingMessages = removeForProcessingMessages(node.ProcessingMessages, 0)
	return m
}
