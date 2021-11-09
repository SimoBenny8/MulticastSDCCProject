package test

import (
	"MulticastSDCCProject/pkg/SQMulticast"
	client2 "MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/testUtil"
	"MulticastSDCCProject/pkg/util"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func count(msgId string) int {

	var c int

	arrayMessage := SQMulticast.GetMessageToBeDeliverQueue()
	for i := range arrayMessage {
		if arrayMessage[i].Id == msgId {
			c += 1
		}
	}
	return c
}

func Assertion(wg *sync.WaitGroup) {

}

func TestOneToMany(t *testing.T) {
	var localErr error
	var connections []*client2.Client
	var wg sync.WaitGroup

	numNode := 3
	messages := []*rpc.Packet{{Message: []byte("message")}, {Message: []byte("in")}, {Message: []byte("order")}}
	connections = make([]*client2.Client, 3)
	msgId := make([]string, 3)
	for i := 0; i < numNode; i++ {
		connections[i] = testUtil.FakeConnect("Node" + strconv.Itoa(i))
	}

	if SQMulticast.SeqPort == nil {
		//scelgo quale dei nodi è il sequencer(randomicamente)
		n := rand.Intn(len(connections))
		SQMulticast.SeqPort = connections[n]
		log.Println("Sequencer is", SQMulticast.SeqPort.Connection.Target())
	}

	go SQMulticast.DeliverSeq(connections)
	//caso invio al sequencer da un nodo generico

	for i := range messages {

		wg.Add(1)
		msgId[i] = SQMulticast.RandSeq(5)
		md := make(map[string]string)
		md[util.TYPEMC] = util.SQMULTICAST
		md[util.TYPENODE] = util.SEQUENCER //a chi arriva
		md[util.MESSAGEID] = msgId[i]

		metad := metadata.New(md)
		ctxNew := metadata.NewOutgoingContext(context.Background(), metad)

		i := i
		go func() {
			_, localErr = SQMulticast.SeqPort.Client.SendPacket(ctxNew, messages[i])
			if localErr != nil {
				t.Errorf("SendPacket failed: %v", localErr)
				return

			}
			wg.Done()
		}()
		wg.Wait()
	}

	//assertion
	time.Sleep(time.Second)
	assert.Equal(t, 9, len(SQMulticast.GetMessageToBeDeliverQueue()))
	for i := range msgId {
		assert.Equal(t, 3, count(msgId[i]))
	}
}
