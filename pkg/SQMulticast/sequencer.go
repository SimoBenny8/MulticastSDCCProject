package SQMulticast

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
	"log"
	"math/rand"
	"sync"
	"time"
)

type MessageSeq struct {
	Message   *rpc.Packet
	Timestamp uint32
	Id        string
}

type Sequencer struct {
	Node           NodeForSq
	LocalTimestamp uint32
	MessageQueue   []MessageSeq
	LocalErr       error
	SeqPort        *client.Client
	Connections    []*client.Client
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	seq     Sequencer
)

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetSequencer() Sequencer {
	return seq
}

func SetSequencer(sequencer Sequencer) {
	seq = sequencer
}

func (s *Sequencer) UpdateTimestamp() {
	s.LocalTimestamp += 1
}

func ReceiveMessageToSequencer(mex *rpc.Packet, id string) {
	var wg sync.Mutex
	seq.UpdateTimestamp()
	log.Println("Timestamp:", seq.LocalTimestamp)
	messageT := &MessageSeq{Message: mex, Timestamp: seq.LocalTimestamp, Id: id}
	wg.Lock()
	seq.addMessageSeq(&wg, messageT)

}

func (s *Sequencer) addMessageSeq(wg *sync.Mutex, mex *MessageSeq) {
	defer wg.Unlock()
	s.MessageQueue = append(s.MessageQueue, *mex)
}

func DeliverSeq(delay int) {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	for {
		//log.Println("esecuzione deliver")
		if len(seq.MessageQueue) > 0 {
			log.Println("sto nella coda")
			message := seq.MessageQueue[0]
			seq.MessageQueue = seq.MessageQueue[1:]
			for i := range seq.Connections {
				wg.Add(1)
				//caso invio al sequencer da un nodo generico
				i := i
				go func() {
					defer wg.Done()
					md := make(map[string]string)
					md[util.TYPEMC] = util.SQMULTICAST
					md[util.TYPENODE] = util.MEMBER //a chi arriva
					md[util.MESSAGEID] = message.Id
					md[util.RECEIVER] = seq.Connections[i].Connection.Target()
					seq.LocalErr = seq.Connections[i].Send(md, message.Message.Message, nil, delay)
					if seq.LocalErr != nil {
						log.Println(seq.LocalErr.Error())
					}
				}()
			}
			wg.Wait()
		}

	}
}
