package impl

import (
	"MulticastSDCCProject/pkg/rpc"
)

var orderedAck []rpc.Packet

/*func init(){
	orderedAck = make([]rpc.Packet,100)
}*/

func GetOrderedAck() []rpc.Packet {
	return orderedAck
}

func AppendOrderedAck(ack rpc.Packet) {
	orderedAck = append(orderedAck, ack)
}

func EmptyOrderedAck() { //svuota l'array
	orderedAck = orderedAck[:0]
}
