package VectorClockMulticast

import (
	"log"
	"sync"
)

type VectorMessages []MessageVectorTimestamp

// Create
var queue VectorMessages

func GetQueue() VectorMessages {
	return queue
}

func init() {
	queue = make(VectorMessages, 0, 100)
}

func AddToQueue(m *MessageVectorTimestamp, wg *sync.Mutex) {
	defer wg.Unlock()
	log.Println("messaggio aggiunto in coda")
	queue = append(queue, *m)
}

func Dequeue() MessageVectorTimestamp {
	log.Println("prendo il primo elemento della coda")
	m := queue[0]
	if len(queue) > 1 {
		queue = queue[1:]
	} else {
		queue = queue[:0]
	}
	return m
}
