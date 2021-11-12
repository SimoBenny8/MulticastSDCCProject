package server

import (
	"MulticastSDCCProject/pkg/MulticastScalarClock/impl"
	"MulticastSDCCProject/pkg/SQMulticast"
	"MulticastSDCCProject/pkg/VectorClockMulticast"
	"MulticastSDCCProject/pkg/rpc"
	"MulticastSDCCProject/pkg/util"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"strconv"
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

func Register(s grpc.ServiceRegistrar) (err error) {
	rpc.RegisterPacketServiceServer(s, NewServer())
	return
}

// RunServer runs gRPC service to publish OPacket service
func RunServer(port uint, gServices ...func(grpc.ServiceRegistrar) error) error {

	netListener, err := getNetListener(port)
	if err == nil {
		log.Println("Succed to listen : ", port)
	}
	// register service
	server := grpc.NewServer()

	for _, grpcService := range gServices {
		err = grpcService(server)
		if err != nil {
			return err
		}
	}

	//svr := NewServer()
	//rpc.RegisterPacketServiceServer(server, svr)
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
			switch values := md.Get(util.TYPEMC); values[0] {
			case util.SCMULTICAST:
				log.Println("case SCMulticast")
				if md.Get(util.ACK)[0] == util.TRUE {
					impl.AppendOrderedAck(message)
				} else if md.Get(util.DELIVER)[0] == util.TRUE {
					log.Println("deliver called for message: " + string(message.Message))
					//util.AppendMessageToBeDeliver(message)
				} else {
					impl.AddingRecevingMex(message)
				}
			case util.SQMULTICAST:
				log.Println("case SQMulticast")
				if md.Get(util.TYPENODE)[0] == util.MEMBER {
					sequencer := SQMulticast.GetSequencer()
					nodes := SQMulticast.GetDeliverNodes()
					messageT := &SQMulticast.MessageSeq{Message: message, Timestamp: sequencer.LocalTimestamp, Id: md.Get(util.MESSAGEID)[0]}
					for i := range nodes {
						if md.Get(util.RECEIVER)[0] == nodes[i].MyConn.Connection.Target() {
							nodes[i].AppendMessageToBeDeliver(messageT)
							log.Println("length queue:" + strconv.Itoa(len(nodes[i].DeliverQueue)))
						}

					}
					log.Println("deliver called for message: " + string(message.Message))
				} else if md.Get(util.TYPENODE)[0] == util.SEQUENCER {
					//arriva al sequencer
					SQMulticast.ReceiveMessageToSequencer(message, md.Get(util.MESSAGEID)[0])
				}
			case util.BMULTICAST:
				log.Println("case BMulticast")
				log.Println("deliver called for message: " + string(message.Message))
			case util.VCMULTICAST:
				log.Println("case VCMulticast")
				if md.Get(util.DELIVER)[0] == util.TRUE {
					log.Println("deliver called for message: " + string(message.Message))
				} else {
					VectorClockMulticast.ReceiveMessage(message)
				}
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
	if !strings.Contains(string(packet.Message), "ack") && impl.GetQueue().Len() > 0 { //TODO: sistemare priorità dei messaggi di deliver
		mexInQueue := impl.Dequeue()
		log.Println("deliver called for message: " + string(mexInQueue.OPacket.Message))
		return DeliverObject{packet}
	} else if strings.Contains(string(packet.Message), "ack") {
		return DeliverObject{}
	}
	return DeliverObject{}
}
