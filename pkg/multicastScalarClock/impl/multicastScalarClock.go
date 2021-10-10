package impl

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/rpc"
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
)

type MessageTimestamp struct {
	Address   uint
	OPacket   rpc.Packet
	Timestamp int32
	Id        string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ProcessingMessage []MessageTimestamp
var respChannel chan []byte

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func SendMessageToAll(message *MessageTimestamp, c []*client.Client, timestamp *int32, wg *sync.WaitGroup) {

	message.Timestamp += 1
	*timestamp += 1

	ProcessingMessage = append(ProcessingMessage, *message)
	if len(ProcessingMessage) > 1 {
		sort.SliceStable(ProcessingMessage, func(i, j int) bool {
			return ProcessingMessage[i].Timestamp < ProcessingMessage[j].Timestamp
		})
	}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network..

	// Encode (send) the value.
	err := enc.Encode(ProcessingMessage[0])
	if len(ProcessingMessage) > 1 {
		ProcessingMessage = ProcessingMessage[1:]
	}
	if err != nil {
		log.Fatal("encode error:", err)
	}
	//respChannel = make(chan []byte,1)
	md := make(map[string]string)
	md["TypeMulticast"] = "SCMulticast"
	for i := range c {
		err = c[i].Send(md, network.Bytes(), respChannel) //grpc
		result := <-respChannel
		log.Println("ack: ", string(result))

		if err != nil {
			log.Fatal("error during send message")
		}
		//if c[i].Connection.Target() != "localhost:"+strconv.FormatUint(uint64(message.Address), 10){
		//	log.Println("destinatario:", c[i].Connection.Target())
		//	SendAck(DecodeMsg(network),timestamp,c[i],nil)
		//}
	}
	//inserire in coda locale
	AddToQueue(DecodeMsg(network))

	wg.Done()

}

func DecodeMsg(network bytes.Buffer) *MessageTimestamp {
	// Decode (receive) the value.
	//log.Println("decodifica del messaggio")
	var err error
	dec := gob.NewDecoder(&network) // Will read from network.
	var m MessageTimestamp
	err = dec.Decode(&m)
	log.Println("messaggio codificato: ", string(m.OPacket.Message))
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return &m

}

func SendAck(mReceived *MessageTimestamp, t *int32, client *client.Client, wg *sync.WaitGroup) {

	if client.Connection.Target() != "localhost:"+strconv.FormatUint(uint64(mReceived.Address), 10) {
		var newTimestamp int32
		if mReceived.Timestamp > *t {
			newTimestamp = mReceived.Timestamp
			*t = newTimestamp
			log.Println("timestamp updated:", newTimestamp)
		} else {
			newTimestamp = *t
			log.Println("timestamp not updated:", newTimestamp)
		}
		log.Println("Original Message: ", string(mReceived.OPacket.Message), "Timestamp: ", mReceived.Timestamp)

		mReceived.OPacket.Message = []byte("ack: " + string(mReceived.OPacket.Message))
		newTimestamp += 1
		*t += 1
		log.Println("timestamp aggiornato per ack: ", newTimestamp)
		ack := &MessageTimestamp{OPacket: mReceived.OPacket, Timestamp: newTimestamp, Id: mReceived.Id}
		var network bytes.Buffer        // Stand-in for a network connection
		enc := gob.NewEncoder(&network) // Will write to network.

		// Encode (send) the value.
		err := enc.Encode(ack)
		if err != nil {
			log.Fatal("encode error:", err)
		}

		err = client.Send(nil, network.Bytes(), nil) //grpc
		if err != nil {
			log.Fatal("error during send ack")
		}
		//TODO: bisogna gestire la deliver del messaggio inviato al processo stesso
		//wg.Done()
	}

}

/* Workflow:
 1) Invio messaggio in multicast da parte di Pi
2) Pj, che riceve il messaggio, lo mette in coda e riordina la coda in base al timestamp (incrementare prima di inviare)
3) Pj invia un messaggio di ack in multicast
4) Pj effettua la deliver se msg è in testa e nessun processo ha messaggio con timestamp minore o uguale in coda */
