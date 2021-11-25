package VectorClockMulticast

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"sync"
)

var (
	Nodes []NodeVC
)

func init() {

	Nodes = make([]NodeVC, 0, 100)

}

func GetNodes() []NodeVC {
	return Nodes
}

//function that return the position of a node in the array
func checkPositionNode(id uint) int {
	nodes := GetNodes()
	for i := range nodes {
		if nodes[i].NodeId == id {
			return i
		}
	}
	return -1
}

func AppendNodes(node NodeVC, wg *sync.Mutex) {
	Nodes = append(Nodes, node)
	wg.Unlock()
}

func (node *NodeVC) AppendDeliverQueue(mex *rpc.Packet) {
	var wg sync.Mutex
	wg.Lock()
	m := DecodeMsg(mex, &wg)
	node.DeliverQueue = append(node.DeliverQueue, *m)
}
