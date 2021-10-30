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
var queue OrderedMessages

func GetQueue() OrderedMessages {
	return queue
}

func init() {
	queue = make(OrderedMessages, 0, 100)
}

func AddToQueue(m *MessageTimestamp) {

	log.Println("messaggio aggiunto in coda")
	queue = append(queue, *m)
	queue = OrderingMessage(queue)
	return
}

func Dequeue() MessageTimestamp {
	log.Println("prendo il primo elemento della coda")
	m := queue[0]
	if len(queue) > 1 {
		queue = queue[1:]
	} else {
		queue = queue[:0]
	}
	return m
}
