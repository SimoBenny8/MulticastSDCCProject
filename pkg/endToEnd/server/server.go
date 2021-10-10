package server

import (
	"MulticastSDCCProject/pkg/SQMulticast"
	"MulticastSDCCProject/pkg/multicastScalarClock/impl"
	"MulticastSDCCProject/pkg/rpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"strings"
	"sync"
)

type Server struct {
	rpc.UnimplementedPacketServiceServer
}

type DeliverObject struct {
	message rpc.Packet
}

func NewServer() *Server {
	return &Server{}
}

// RunServer runs gRPC service to publish OPacket service
func RunServer(port uint) error {

	netListener, err := getNetListener(port)
	if err == nil {
		log.Println("Succed to listen : ", port)
	}
	// register service
	server := grpc.NewServer()
	svr := NewServer()
	rpc.RegisterPacketServiceServer(server, svr)
	// start the server
	err = server.Serve(netListener)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	} else {
		log.Println("server connected")
	}
	return err
}

func RunServerWithWaitGroup(port uint, wg *sync.WaitGroup) error {

	netListener, err := getNetListener(port)
	if err == nil {
		log.Println("Succed to listen : ", port)
	}
	// register service
	server := grpc.NewServer()
	svr := NewServer()
	rpc.RegisterPacketServiceServer(server, svr)
	// start the server
	err = server.Serve(netListener)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	} else {
		log.Println("server connected")
	}
	wg.Done()
	return err
}

func getNetListener(port uint) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	return lis, err
}

//implementation of service message
func (s *Server) SendPacket(ctx context.Context, message *rpc.Packet) (*rpc.ResponsePacket, error) {
	log.Println("Received Message: ", string(message.Message))
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md) > 0 {
			switch values := md.Get(SQMulticast.TYPEMC); values[0] {
			case SQMulticast.SCMULTICAST:
				log.Println("case SCMulticast")
			case SQMulticast.SQMULTICAST:
				log.Println("case SQMulticast")
				if md.Get(SQMulticast.TYPENODE)[0] == SQMulticast.MEMBER {
					if SQMulticast.LocalTimestamp == 0 {
						SQMulticast.InitTimestamp()
						SQMulticast.InitMessageQueue()
						SQMulticast.UpdateTimestamp()
						messageT := &SQMulticast.MessageT{Message: *message, Timestamp: SQMulticast.LocalTimestamp}
						SQMulticast.MessageQueue = append(SQMulticast.MessageQueue, *messageT)
						go SQMulticast.DeliverSeq()
					} else {
						SQMulticast.UpdateTimestamp()
						messageT := &SQMulticast.MessageT{Message: *message, Timestamp: SQMulticast.LocalTimestamp}
						SQMulticast.MessageQueue = append(SQMulticast.MessageQueue, *messageT)
						go SQMulticast.DeliverSeq()
					}
				} else if md.Get(SQMulticast.TYPENODE)[0] == SQMulticast.SEQUENCER { //arriva dal sequencer
					log.Println("deliver called for message: " + string(message.Message))
				}
			case SQMulticast.BMULTICAST:
				log.Println("case BMulticast")
				log.Println("deliver called for message: " + string(message.Message))
			case SQMulticast.VCMULTICAST:
				log.Println("case VCMulticast")
			default:
				panic("unrecognized value")

			}
		} else {
			log.Fatal("metadata not arrived")
		}
	}
	return &rpc.ResponsePacket{}, nil
}

func DeliverMulticast(packet rpc.Packet) DeliverObject {
	if !strings.Contains(string(packet.Message), "ack") && impl.Queue.Len() > 0 { //TODO: sistemare priorità dei messaggi di deliver
		mexInQueue := impl.Dequeue()
		log.Println("deliver called for message: " + string(mexInQueue.OPacket.Message))
		return DeliverObject{packet}
	} else if strings.Contains(string(packet.Message), "ack") {
		return DeliverObject{}
	}
	return DeliverObject{}
}
