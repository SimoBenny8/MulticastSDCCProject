package impl

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"sync"
)

type AckNode struct {
	id          string
	connections []*client.Client
	mutex       sync.Mutex
}

var Node *AckNode

var orderedAck []MessageTimestamp

func init() {
	orderedAck = make([]MessageTimestamp, 0, 100)
	Node = &AckNode{
		id:          "",
		connections: connections,
		mutex:       sync.Mutex{},
	}
}

func AppendOrderedAck(ack *rpc.Packet) {
	m := DecodeMsg(ack)
	orderedAck = append(orderedAck, *m)
	return
}

func EmptyOrderedAck() { //svuota l'array
	orderedAck = orderedAck[:0]
	return
}
