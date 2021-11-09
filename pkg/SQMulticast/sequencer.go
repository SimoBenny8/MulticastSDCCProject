package SQMulticast

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"sync"
	"time"
)

type MessageT struct {
	Message   rpc.Packet
	Timestamp uint32
	Id        string
}

var (
	LocalTimestamp uint32
	MessageQueue   []MessageT
	letters        = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	LocalErr       error
	SeqPort        *client.Client
)

func init() {
	MessageQueue = make([]MessageT, 0, 100)
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

func DeliverSeq(Connections []*client.Client) {
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
					md[util.TYPEMC] = util.SQMULTICAST
					md[util.TYPENODE] = util.MEMBER //a chi arriva
					md[util.MESSAGEID] = message.Id

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
