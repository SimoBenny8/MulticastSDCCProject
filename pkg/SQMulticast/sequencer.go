package SQMulticast

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"sync"
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
	Id        string
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

func UpdateTimestamp() {
	LocalTimestamp += 1
}

func DeliverSeq() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	for {
		if len(MessageQueue) > 0 {
			message := MessageQueue[0]
			MessageQueue = MessageQueue[1:]
			for i := range Connections {
				wg.Add(1)
				//caso invio al sequencer da un nodo generico
				i := i
				go func() {
					defer wg.Done()
					md := make(map[string]string)
					md[TYPEMC] = SQMULTICAST
					md[TYPENODE] = MEMBER //a chi arriva
					md[MESSAGEID] = message.Id
					delay := rand.Intn(10700) + 1000
					//log.Println("Delay: ",delay," milliseconds")
					time.Sleep(time.Duration(delay))
					metaData := metadata.New(md)
					ctx := metadata.NewOutgoingContext(context.Background(), metaData)
					_, LocalErr = Connections[i].Client.SendPacket(ctx, &message.Message)
					if LocalErr != nil {
						log.Println(LocalErr.Error())
					}
				}()
			}
			wg.Wait()
		}

	}
}
