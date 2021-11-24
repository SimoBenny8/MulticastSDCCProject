package MulticastScalarClock

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
)

type OtherTimestamp struct {
	id                 string
	otherNodeTimestamp int
}

type NodeSC struct {
	NodeId             uint
	Connections        []*client.Client
	ProcessingMessages OrderedMessages
	MyConn             *client.Client
	ReceivedMessage    OrderedMessages
	Timestamp          int
	OtherTs            []OtherTimestamp
	OrderedAck         []MessageTimestamp
	DeliverQueue       OrderedMessages
}

var (
	Nodes []NodeSC
)

func init() {

	Nodes = make([]NodeSC, 0, 100)

}

func GetNodes() []NodeSC {
	return Nodes
}

func (node *NodeSC) AppendDeliverMessages(mex *rpc.Packet) {
	m := DecodeMsg(mex)
	//pos := checkPositionNode(nodeId)
	node.DeliverQueue = append(node.DeliverQueue, *m)
}

func (node *NodeSC) AppendOrderedAck(ack *rpc.Packet) {
	m := DecodeMsg(ack)
	//pos := checkPositionNode(nodeId)
	node.OrderedAck = append(node.OrderedAck, *m)
	node.OtherTs = append(node.OtherTs, OtherTimestamp{
		id:                 m.Id,
		otherNodeTimestamp: m.FirstTsInQueue,
	})
	return
}

func checkPositionNode(id uint) int {
	for i := range Nodes {
		if Nodes[i].NodeId == id {
			return i
		}
	}
	return -1
}

func removeForReceivedMessage(slice OrderedMessages, s int) OrderedMessages {
	return append(slice[:s], slice[s+1:]...)
}

func removeForProcessingMessages(slice OrderedMessages, s int) OrderedMessages {
	return append(slice[:s], slice[s+1:]...)
}

func removeForMessageTimestamp(slice []MessageTimestamp, s int) []MessageTimestamp {
	return append(slice[:s], slice[s+1:]...)
}

func removeForOtherTimestamps(slice []OtherTimestamp, s int) []OtherTimestamp {
	return append(slice[:s], slice[s+1:]...)
}

func EmptyOtherTimestamp(idMex string, nodeId uint) {
	pos := checkPositionNode(nodeId)
	for i := range Nodes[pos].OtherTs {
		if Nodes[pos].OtherTs[i].id == idMex {
			Nodes[pos].OtherTs = removeForOtherTimestamps(Nodes[pos].OtherTs, i)
		}
	}

}

func AppendNodes(node NodeSC) {
	Nodes = append(Nodes, node)

}

func EmptyOrderedAck(idMex string, nodeId uint) { //svuota l'array
	pos := checkPositionNode(nodeId)
	for i := range Nodes[pos].OrderedAck {
		if Nodes[pos].OrderedAck[i].Id == idMex {
			Nodes[pos].OrderedAck = removeForMessageTimestamp(Nodes[pos].OrderedAck, i)
		}
	}
	return
}
