package SQMulticast

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"time"
)

const (
	SEQUENCER   = "SequencerNode"
	SQMULTICAST = "SQMulticast"
	SCMULTICAST = "SCMulticast"
	BMULTICAST  = "BMulticast"
	VCMULTICAST = "VCMulticast"
	TYPEMC      = "TypeMulticast"
	MEMBER      = "MemberNode"
	TYPENODE    = "TypeNode"
	MESSAGEID   = "MessageId"
)

type MessageT struct {
	Message   rpc.Packet
	Timestamp uint32
}

var LocalTimestamp uint32
var MessageQueue []MessageT
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var Connections []*client.Client
var LocalErr error
var SeqPort *client.Client

func InitMessageQueue() {
	MessageQueue = make([]MessageT, 100)
}

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func InitTimestamp() {
	LocalTimestamp = 0
}

func UpdateTimestamp() {
	LocalTimestamp += 1
}

/*func ReceiveMessage(messageCh chan rpc.Packet){
	select{
		case <- messageCh:
			LocalTimestamp += 1
			messageWithT := &MessageT{Message: <-messageCh,Timestamp: LocalTimestamp}
			MessageQueue = append(MessageQueue,*messageWithT)
			if len(MessageQueue) > 1{
				sort.Slice(MessageQueue, func(i, j int) bool {
					return MessageQueue[i].Timestamp < MessageQueue[j].Timestamp
				})
			}


	}
}*/

func DeliverSeq() {
	//mettere waitgroup
	if len(MessageQueue) > 0 {
		message := MessageQueue[0]
		if len(MessageQueue) > 1 {
			MessageQueue = MessageQueue[1:]
		}
		for i := range Connections {
			if Connections[i].Connection.Target() != SeqPort.Connection.Target() {
				//caso invio al sequencer da un nodo generico
				md := make(map[string]string)
				md[TYPEMC] = SQMULTICAST
				md[TYPENODE] = SEQUENCER //da chi arriva
				go time.Sleep(time.Duration(rand.Intn(3700) + 5))
				metaData := metadata.New(md)
				ctx := metadata.NewOutgoingContext(context.Background(), metaData)
				_, LocalErr = Connections[i].Client.SendPacket(ctx, &message.Message)
				if LocalErr != nil {
					log.Println(LocalErr.Error())
				}
			} else {
				continue
			}
		}
	}
}
