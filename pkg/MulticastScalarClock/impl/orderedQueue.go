package impl

import (
	"log"
	"sort"
)

type OrderedMessages []MessageTimestamp

func (a OrderedMessages) Len() int           { return len(a) }
func (a OrderedMessages) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }
func (a OrderedMessages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//func che ordina i messaggi in coda
func (node *NodeSC) OrderingMessage(messages []MessageTimestamp) OrderedMessages {
	log.Println("riordino i messaggi")
	sort.Sort(node.DeliverQueue)
	return messages
}

// Create

func init() {
	//queue = make(OrderedMessages, 0, 100)
}

func (node *NodeSC) AddToQueue(m *MessageTimestamp) {

	log.Println("messaggio aggiunto in coda")
	node.DeliverQueue = append(node.DeliverQueue, *m)
	node.DeliverQueue = node.OrderingMessage(node.DeliverQueue)
	return
}

func (node *NodeSC) Dequeue() MessageTimestamp {
	log.Println("prendo il primo elemento della coda")
	m := node.DeliverQueue[0]
	if len(node.DeliverQueue) > 1 {
		node.DeliverQueue = node.DeliverQueue[1:]
	} else {
		node.DeliverQueue = node.DeliverQueue[:0]
	}
	return m
}
