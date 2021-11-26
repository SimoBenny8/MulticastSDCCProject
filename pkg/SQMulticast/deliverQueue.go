package SQMulticast

import (
	"sort"
)

var (
	Nodes []NodeForSq
)

type OrderedMessages []*MessageSeq

func (a OrderedMessages) Len() int           { return len(a) }
func (a OrderedMessages) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }
func (a OrderedMessages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//Ordering messages queue
func (node *NodeForSq) OrderingMessage(slice OrderedMessages) {
	sort.Sort(slice)
	return
}

func init() {
	Nodes = make([]NodeForSq, 0, 100)
}

func GetDeliverNodes() []NodeForSq {
	return Nodes
}

func (node *NodeForSq) GetMessageToBeDeliverQueue() []*MessageSeq {
	return node.DeliverQueue
}

//add a message to deliver queue
func (node *NodeForSq) AppendMessageDelivered(m *MessageSeq) {

	node.DeliverQueue = append(node.DeliverQueue, m)
	node.OrderingMessage(node.DeliverQueue)
}

//add a message to processing queue
func (node *NodeForSq) AppendMessageToBeDeliver(m *MessageSeq) {

	node.ProcessingQueue = append(node.ProcessingQueue, m)
	node.OrderingMessage(node.ProcessingQueue)
}

func AppendNodes(node NodeForSq) {
	Nodes = append(Nodes, node)

}
