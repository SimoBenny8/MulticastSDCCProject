package SQMulticast

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"sort"
)

var (
	Nodes []NodeForSq
)

type NodeForSq struct {
	NodeId       uint
	Connections  []*client.Client
	DeliverQueue OrderedMessages
	MyConn       *client.Client
}

type OrderedMessages []*MessageSeq

func (a OrderedMessages) Len() int           { return len(a) }
func (a OrderedMessages) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }
func (a OrderedMessages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//Ordering messages queue
func (node *NodeForSq) OrderingMessage() {
	sort.Sort(node.DeliverQueue)
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
func (node *NodeForSq) AppendMessageToBeDeliver(m *MessageSeq) {

	node.DeliverQueue = append(node.DeliverQueue, m)
	node.OrderingMessage()
}

func AppendNodes(node NodeForSq) {
	Nodes = append(Nodes, node)

}
