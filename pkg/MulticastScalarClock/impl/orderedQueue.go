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
func OrderingMessage(messages []MessageTimestamp) OrderedMessages { //da testare
	log.Println("riordino i messaggi")
	sort.Sort(OrderedMessages(messages))
	return messages
}

// Create
var Queue OrderedMessages

func InitQueue() {
	Queue = make(OrderedMessages, 50)
}

func AddToQueue(m *MessageTimestamp) OrderedMessages {

	log.Println("messaggio aggiunto in coda")
	Queue = append(Queue, *m)
	Queue = OrderingMessage(Queue)
	return Queue
}

func Dequeue() MessageTimestamp {
	log.Println("prendo il primo elemento della coda")
	m := Queue[0]
	Queue = Queue[1:]
	return m
}
