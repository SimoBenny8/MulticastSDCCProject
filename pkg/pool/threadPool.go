package pool

import (
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/MulticastScalarClock"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/SQMulticast"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/VectorClockMulticast"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"log"
	"sync"
)

type ThreadPool struct {
	Wg      sync.Mutex
	Message chan *rpc.Packet
}

var (
	Pool ThreadPool
)

func (t *ThreadPool) InitThreadPool(connections []*client.Client, numTh int, multicastType string, respChannel chan []byte, port uint, nodeId uint, delay int) {
	t.Message = make(chan *rpc.Packet, 30)

	Pool.Wg.Lock()

	for i := 0; i < numTh; i++ {
		go getMessages(Pool.Message, multicastType, connections, respChannel, port, nodeId, delay)
	}

	Pool.Wg.Unlock()

}

func getMessages(chanMex chan *rpc.Packet, multicastType string, connections []*client.Client, respChannel chan []byte, port uint, nodeId uint, delay int) {
	var localErr error
	seq := SQMulticast.GetSequencer()
	for {
		select {
		case mex := <-chanMex:
			if multicastType == util.SQMULTICAST {
				for i := range connections {
					if connections[i].Connection.Target() == seq.SeqPort.Connection.Target() {
						//caso invio al sequencer da un nodo generico
						md := make(map[string]string)
						md[util.TYPEMC] = util.SQMULTICAST
						md[util.TYPENODE] = util.SEQUENCER //a chi arriva
						md[util.MESSAGEID] = SQMulticast.RandSeq(5)
						localErr = connections[i].Send(md, []byte(""), mex.Message, nil, delay)
					}
				}
				if localErr != nil {
					log.Println("Error in sending to node")
				}
			} else if multicastType == util.SCMULTICAST {
				message := &MulticastScalarClock.MessageTimestamp{Address: port, OPacket: *mex, Timestamp: MulticastScalarClock.GetTimestamp(nodeId), Id: MulticastScalarClock.RandSeq(5)}
				MulticastScalarClock.SendMessageToAll(message, nodeId, delay)

			} else if multicastType == util.VCMULTICAST {
				VectorClockMulticast.SendMessageToAll(mex, port, nodeId, delay)

			} else if multicastType == util.BMULTICAST {
				for i := range connections {
					md := make(map[string]string)
					md[util.TYPEMC] = util.BMULTICAST
					localErr = connections[i].Send(md, mex.Header, mex.Message, respChannel, delay)
					result := <-respChannel
					fmt.Println(string(result)) //problema ack implosion
				}
				if localErr != nil {
					log.Println("Error in sending to node")
				}
			}

		}

	}
}
