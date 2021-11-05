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

type OtherTimestamp struct {
	id                 string
	otherNodeTimestamp int
}

var otherTs []OtherTimestamp

var Node *AckNode

var orderedAck []MessageTimestamp

func init() {
	orderedAck = make([]MessageTimestamp, 0, 100)
	otherTs = make([]OtherTimestamp, 0, 100)
	Node = &AckNode{
		id:          "",
		connections: connections,
		mutex:       sync.Mutex{},
	}
}

func AppendOrderedAck(ack *rpc.Packet) {
	m := DecodeMsg(ack)
	orderedAck = append(orderedAck, *m)
	otherTs = append(otherTs, OtherTimestamp{
		id:                 m.Id,
		otherNodeTimestamp: m.FirstTsInQueue,
	})
	return
}

func removeForMessageTimestamp(slice []MessageTimestamp, s int) []MessageTimestamp {
	return append(slice[:s], slice[s+1:]...)
}

func removeForOtherTimestamps(slice []OtherTimestamp, s int) []OtherTimestamp {
	return append(slice[:s], slice[s+1:]...)
}

func EmptyOtherTimestamp(idMex string) {
	for i := range otherTs {
		if otherTs[i].id == idMex {
			removeForOtherTimestamps(otherTs, i)
		}
	}

	return

}

func EmptyOrderedAck(idMex string) { //svuota l'array

	for i := range orderedAck {
		if orderedAck[i].Id == idMex {
			removeForMessageTimestamp(orderedAck, i)
		}
	}
	return
}
