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

//func che ordina i messaggi in coda
func (node *NodeForSq) OrderingMessage() {
	//log.Println("riordino i messaggi")
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

func (node *NodeForSq) AppendMessageToBeDeliver(m *MessageSeq) {

	node.DeliverQueue = append(node.DeliverQueue, m)
	node.OrderingMessage()
}

func AppendNodes(node NodeForSq) {
	Nodes = append(Nodes, node)

}
