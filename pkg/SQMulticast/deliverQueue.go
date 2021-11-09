package SQMulticast

import (
	"sort"
)

var (
	message []*MessageT
)

type OrderedMessages []*MessageT

func (a OrderedMessages) Len() int           { return len(a) }
func (a OrderedMessages) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }
func (a OrderedMessages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//func che ordina i messaggi in coda
func OrderingMessage(messages OrderedMessages) OrderedMessages {
	//log.Println("riordino i messaggi")
	sort.Sort(messages)
	return messages
}

func init() {
	message = make([]*MessageT, 0, 100)
}

func GetMessageToBeDeliverQueue() []*MessageT {
	return message
}

func AppendMessageToBeDeliver(m *MessageT) {
	message = append(message, m)
	OrderingMessage(message)
}
