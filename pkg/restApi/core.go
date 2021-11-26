package restApi

import (
	"errors"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/MulticastScalarClock"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/SQMulticast"
	client1 "github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/proto"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/VectorClockMulticast"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/client"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/pool"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MulticastGroup struct {
	ClientId  string
	Group     *MulticastInfo
	groupMu   sync.RWMutex
	Messages  []Message
	MessageMu sync.RWMutex
}

type Message struct {
	MessageHeader map[string]string `json:"MessageHeader"`
	Payload       []byte            `json:"Payload"`
}

type MessageRequest struct {
	MulticastId string `json:"multicast_id"`
	message     Message
}

type MulticastInfo struct {
	MulticastId      string            `json:"multicast_id"`
	MulticastType    string            `json:"multicast_type"`
	ReceivedMessages int               `json:"received_messages"`
	Status           string            `json:"status"`
	Members          map[string]Member `json:"members"`
}

type Member struct {
	MemberId string `json:"member_id"`
	Address  string `json:"address"`
	Ready    bool   `json:"ready"`
}

type MulticastReq struct {
	MulticastId   string `json:"multicast_id"`
	MulticastType string `json:"multicast_type"`
}
type MulticastId struct {
	MulticastId string `json:"multicast_id"`
}

type GroupConfig struct {
	MulticastType string `json:"multicast_type"`
}

var (
	RegClient       proto.RegistryClient
	GMu             sync.RWMutex
	Delay           uint
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
	NumThread       int
)

func init() {
	MulticastGroups = make(map[string]*MulticastGroup)
}

func Run(grpcP uint, restPort uint, registryAddr, relativePath string, numThreads int, dl uint, debug bool) error {
	GrpcPort = grpcP
	Delay = dl
	NumThread = numThreads

	var err error
	RegClient, err = client1.Connect(registryAddr)

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if err != nil {
		return err
	}

	r := routes{
		router: gin.Default(),
	}

	v1 := r.router.Group(relativePath)
	r.addGroups(v1)
	err = r.router.Run(fmt.Sprintf(":%d", restPort))
	return err
}

func (r routes) addGroups(rg *gin.RouterGroup) {
	groups := rg.Group("/groups")
	groups.POST("/", addGroup)
	groups.PUT("/:mId", startGroup)
	groups.GET("/:mId", getInfoGroup)
	groups.DELETE("/:mId", closeGroup)
	groups.POST("/messages/:mId", sendMessage)
	groups.GET("/messages/:mId", getMessages)
	groups.GET("/deliverQueue/:mId", getDeliverQueue)
}

func InitGroup(info *proto.Group, group *MulticastGroup, port uint) {

	var myConn *client.Client

	log.Println("Waiting for the group to be ready")

	update(info, group)
	groupInfo, err := ChangeStatus(info, group, proto.Status_OPENING)
	if err != nil {
		return
	}

	log.Println("Group ready, initializing multicast")

	var members []string

	for memberId, member := range group.Group.Members {
		if memberId != group.ClientId {
			members = append(members, member.Address)
		}
	}
	members = append(members, group.ClientId)

	connections := make([]*client.Client, len(members))

	for i, member := range members {
		log.Println("Connecting with", members[i])
		connections[i] = client.Connect(member)
		if strings.Contains(connections[i].Connection.Target(), strconv.Itoa(int(port))) {
			myConn = connections[i]
		}
	}

	// Initializing  data structures
	err = initGroupCommunication(info, port, connections, myConn)

	if err != nil {
		return
	}

	groupInfo, err = RegClient.Ready(context.Background(), &proto.RequestData{
		MulticastId: group.Group.MulticastId,
		ClientId:    group.ClientId,
	})
	if err != nil {
		return
	}

	log.Println("Waiting for the other nodes")
	if groupInfo.MulticastType.String() == util.VCMULTICAST {
		nodes := VectorClockMulticast.GetNodes()
		for i := range nodes {
			if nodes[i].MyConn == myConn {
				go VectorClockMulticast.Deliver(nodes[i].NodeId, int(Delay))
			}
		}
	}
	if groupInfo.MulticastType.String() == util.SQMULTICAST {
		go SQMulticast.DeliverSeq(int(Delay))
		nodes := SQMulticast.GetDeliverNodes()
		for i := range nodes {
			if nodes[i].MyConn == myConn {
				go SQMulticast.DeliverMsg(int(Delay), nodes[i].NodeId)
			}
		}
	}
	if groupInfo.MulticastType.String() == util.SCMULTICAST {
		nodes := MulticastScalarClock.GetNodes()
		for i := range nodes {
			if nodes[i].MyConn == myConn {
				go MulticastScalarClock.Deliver(len(connections), nodes[i].NodeId, int(Delay))
			}
		}
	}
	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = ChangeStatus(groupInfo, group, proto.Status_STARTING)

	if err != nil {
		return
	}

	log.Println("Ready to multicast")

}

//Start communication
func initGroupCommunication(groupInfo *proto.Group, port uint, connections []*client.Client, myConn *client.Client) error {

	if groupInfo.MulticastType.String() == util.BMULTICAST {
		log.Println("STARTING BMULTICAST")
		respChannel := make(chan []byte, 1)
		pool.Pool.InitThreadPool(connections, NumThread, util.BMULTICAST, respChannel, port, 0, int(Delay))

	}
	if groupInfo.MulticastType.String() == util.SQMULTICAST {
		log.Println("STARTING SQMULTICAST")
		nodeID, _ := strconv.Atoi(groupInfo.MulticastId)
		n := rand.Intn(len(connections))
		node := new(SQMulticast.NodeForSq)
		node.NodeId = uint(nodeID)
		node.Connections = connections
		node.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
		node.MyConn = myConn
		node.LocalTimestamp = 0
		node.ProcessingQueue = make([]*SQMulticast.MessageSeq, 0, 100)

		SQMulticast.AppendNodes(*node)
		seq := new(SQMulticast.Sequencer)
		seq.Node = *node
		seq.SeqPort = connections[n]
		seq.Connections = connections
		seq.LocalTimestamp = 0
		seq.MessageQueue = make([]SQMulticast.MessageSeq, 0, 100)
		SQMulticast.SetSequencer(*seq)
		pool.Pool.InitThreadPool(connections, NumThread, util.SQMULTICAST, nil, port, node.NodeId, int(Delay))
		log.Println("The sequencer nodes is at port", seq.SeqPort.Connection.Target())
	}

	if groupInfo.MulticastType.String() == util.SCMULTICAST {
		log.Println("STARTING SCMULTICAST")
		node := new(MulticastScalarClock.NodeSC)
		node.NodeId = uint(rand.Intn(5000))
		node.Connections = connections
		node.ProcessingMessages = make(MulticastScalarClock.OrderedMessages, 0, 100)
		node.MyConn = myConn
		node.ReceivedMessage = make(MulticastScalarClock.OrderedMessages, 0, 100)
		node.Timestamp = 0
		node.OrderedAck = make(MulticastScalarClock.OrderedMessages, 0, 100)
		node.OtherTs = make([]MulticastScalarClock.OtherTimestamp, 0, 100)
		node.DeliverQueue = make(MulticastScalarClock.OrderedMessages, 0, 100)

		MulticastScalarClock.AppendNodes(*node)

		pool.Pool.InitThreadPool(connections, NumThread, util.SCMULTICAST, nil, port, node.NodeId, int(Delay))
		go MulticastScalarClock.Receive(node.NodeId, node.NodeId, int(Delay))
	}
	if groupInfo.MulticastType.String() == util.VCMULTICAST {
		log.Println("STARTING VCMULTICAST")
		var wg sync.Mutex
		var myNode int32
		for i := range connections {
			if strings.Contains(connections[i].Connection.Target(), strconv.Itoa(int(port))) {
				myNode = int32(i)
				log.Println("my index: ", myNode)
			}
		}

		node := new(VectorClockMulticast.NodeVC)
		wg.Lock()
		node.InitLocalTimestamp(&wg, len(connections))
		node.NodeId = uint(rand.Intn(5000))
		node.Connections = connections
		node.DeliverQueue = make(VectorClockMulticast.VectorMessages, 0, 100)
		node.MyConn = myConn
		node.MyNode = myNode

		wg.Lock()
		VectorClockMulticast.AppendNodes(*node, &wg)
		pool.Pool.InitThreadPool(connections, NumThread, util.VCMULTICAST, nil, port, node.NodeId, int(Delay))

	}
	return nil

}

//Change status of the member
func ChangeStatus(groupInfo *proto.Group, multicastGroup *MulticastGroup, status proto.Status) (*proto.Group, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 1)
		groupInfo, err = RegClient.GetStatus(context.Background(), &proto.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}
		update(groupInfo, multicastGroup)
	}

	if groupInfo.Status == proto.Status_CLOSED || groupInfo.Status == proto.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}

	return groupInfo, nil
}

//update status of member
func update(groupInfo *proto.Group, multicastGroup *MulticastGroup) {
	multicastGroup.groupMu.Lock()
	defer multicastGroup.groupMu.Unlock()

	multicastGroup.Group.Status = proto.Status_name[int32(groupInfo.Status)]

	for clientId, member := range groupInfo.Members {
		m, ok := multicastGroup.Group.Members[clientId]

		if !ok {
			m = Member{
				MemberId: member.Id,
				Address:  member.Address,
				Ready:    member.Ready,
			}

			multicastGroup.Group.Members[clientId] = m
		}

		if m.Ready != member.Ready {
			m.Ready = member.Ready
			multicastGroup.Group.Members[clientId] = m
		}

	}
}
