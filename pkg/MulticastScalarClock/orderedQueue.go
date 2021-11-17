package MulticastScalarClock

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
	sort.Sort(node.ProcessingMessages)
	return messages
}

// Create

func (node *NodeSC) AddToQueue(m *MessageTimestamp) {

	log.Println("messaggio aggiunto in coda")
	node.ProcessingMessages = append(node.ProcessingMessages, *m)
	node.ProcessingMessages = node.OrderingMessage(node.ProcessingMessages)
	return
}

func (node *NodeSC) Dequeue() MessageTimestamp {
	log.Println("prendo il primo elemento della coda")
	m := node.ProcessingMessages[0]
	if len(node.ProcessingMessages) > 1 {
		node.ProcessingMessages = node.ProcessingMessages[1:]
	} else {
		node.ProcessingMessages = node.ProcessingMessages[:0]
	}
	return m
}
