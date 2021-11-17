package MulticastScalarClock

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
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
	ReceivedMessage    []*rpc.Packet
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

func GetDeliverNodes() []NodeSC {
	return Nodes
}

func AppendDeliverMessages(mex *rpc.Packet, nodeId uint) {
	m := DecodeMsg(mex)
	pos := checkPositionNode(nodeId)
	Nodes[pos].DeliverQueue = append(Nodes[pos].DeliverQueue, *m)
}

func AppendOrderedAck(ack *rpc.Packet, nodeId uint) {
	m := DecodeMsg(ack)
	pos := checkPositionNode(nodeId)
	Nodes[pos].OrderedAck = append(Nodes[pos].OrderedAck, *m)
	Nodes[pos].OtherTs = append(Nodes[pos].OtherTs, OtherTimestamp{
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
	return 0
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
			removeForOtherTimestamps(Nodes[pos].OtherTs, i)
		}
	}

	return

}

func AppendNodes(node NodeSC) {
	Nodes = append(Nodes, node)

}

func EmptyOrderedAck(idMex string, nodeId uint) { //svuota l'array
	pos := checkPositionNode(nodeId)
	for i := range Nodes[pos].OrderedAck {
		if Nodes[pos].OrderedAck[i].Id == idMex {
			removeForMessageTimestamp(Nodes[pos].OrderedAck, i)
		}
	}
	return
}
