package SQMulticast

import (
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"log"
	"math/rand"
	"sync"
	"time"
)

type NodeForSq struct {
	NodeId          uint
	Connections     []*client.Client
	DeliverQueue    OrderedMessages
	MyConn          *client.Client
	LocalTimestamp  int
	ProcessingQueue OrderedMessages
}

//Function that return the array position of the node
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

//Get a message from queue
func (node *NodeForSq) Dequeue() MessageSeq {
	log.Println("Get the first message of node:", node.NodeId)
	m := node.ProcessingQueue[0]
	node.ProcessingQueue = removeForReceivedMessage(node.ProcessingQueue, 0)
	return *m
}

func DeliverMsg(delay int, nodeId uint) {
	pos := checkPositionNode(nodeId)
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	for {
		if len(Nodes[pos].ProcessingQueue) > 0 && int(Nodes[pos].ProcessingQueue[0].Timestamp) <= Nodes[pos].LocalTimestamp {
			message := Nodes[pos].Dequeue()
			wg.Add(1)
			go func() {
				var localErr error
				defer wg.Done()
				md := make(map[string]string)
				md[util.TYPEMC] = util.SQMULTICAST
				md[util.TYPENODE] = util.DELIVER
				md[util.MESSAGEID] = message.Id
				md[util.RECEIVER] = Nodes[pos].MyConn.Connection.Target()
				localErr = Nodes[pos].MyConn.Send(md, []byte(""), message.Message.Message, nil, delay)
				if localErr != nil {
					log.Println(localErr.Error())
				}
			}()
			wg.Wait()
		}

	}

}
