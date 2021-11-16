package VectorClockMulticast

import (
	"MulticastSDCCProject/pkg/rpc"
	"sync"
)

var (
	Nodes []NodeVC
)

func init() {

	Nodes = make([]NodeVC, 0, 100)

}

func checkPositionNode(id uint) int {
	for i := range Nodes {
		if Nodes[i].NodeId == id {
			return i
		}
	}
	return 0
}

func AppendNodes(node NodeVC) {
	Nodes = append(Nodes, node)

}

func AppendDeliverQueue(mex *rpc.Packet, nodeId uint) {
	var wg sync.Mutex
	wg.Lock()
	m := DecodeMsg(mex, &wg)
	pos := checkPositionNode(nodeId)
	Nodes[pos].DeliverQueue = append(Nodes[pos].DeliverQueue, *m)
}
