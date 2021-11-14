package restApi

import (
	"MulticastSDCCProject/pkg/MulticastScalarClock/impl"
	"MulticastSDCCProject/pkg/SQMulticast"
	"MulticastSDCCProject/pkg/ServiceRegistry/ServiceProto"
	client1 "MulticastSDCCProject/pkg/ServiceRegistry/client"
	"MulticastSDCCProject/pkg/VectorClockMulticast"
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/pool"
	"MulticastSDCCProject/pkg/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type MulticastGroup struct {
	clientId  string
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
	RegClient       ServiceProto.RegistryClient
	GMu             sync.RWMutex
	Delay           uint
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
	Application     bool
)

func init() {
	MulticastGroups = make(map[string]*MulticastGroup)
}

func Run(grpcP uint, restPort uint, registryAddr, relativePath string, numThreads int, dl uint, debug bool) error {
	GrpcPort = grpcP
	Delay = dl

	var err error
	Application = true
	//utils2.Vectorclock = utils2.NewVectorClock(4)

	RegClient, err = client1.Connect(registryAddr)

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
}

func InitGroup(info *ServiceProto.Group, group *MulticastGroup) {

	log.Println("Waiting for the group to be ready")

	update(info, group)
	groupInfo, err := ChangeStatus(info, group, ServiceProto.Status_OPENING)
	if err != nil {
		return
	}

	log.Println("Group ready, initializing multicast")
	// Initializing  data structures
	err = initGroupCommunication(info, group)

	if err != nil {
		return
	}

	groupInfo, err = RegClient.Ready(context.Background(), &ServiceProto.RequestData{
		MulticastId: group.Group.MulticastId,
		ClientId:    group.clientId,
	})
	if err != nil {
		return
	}

	log.Println("Waiting for the other nodes")
	if groupInfo.MulticastType.String() == util.VCMULTICAST {
		go VectorClockMulticast.Deliver()
	}
	if groupInfo.MulticastType.String() == util.SQMULTICAST {
		go SQMulticast.DeliverSeq()
	}
	if groupInfo.MulticastType.String() == util.SCMULTICAST {
		go impl.Deliver()
	}
	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = ChangeStatus(groupInfo, group, ServiceProto.Status_STARTING)

	if err != nil {
		return
	}

	log.Println("Ready to multicast")

}

//Start communication
func initGroupCommunication(groupInfo *ServiceProto.Group, group *MulticastGroup) error {
	var members []string

	for memberId, member := range group.Group.Members {
		if memberId != group.clientId {
			members = append(members, member.Address)
		}
	}
	members = append(members, group.clientId)

	connections := make([]*client.Client, len(members))

	for i := range members {
		connections[i] = client.Connect("localhost:" + members[i])
	}

	if groupInfo.MulticastType.String() == "BMULTICAST" {
		log.Println("STARTING BMULTICAST COMMUNICATION")
		respChannel := make(chan []byte, 1)
		pool.Pool.InitThreadPool(connections, 5, util.BMULTICAST, respChannel, *port)

	}
	if groupInfo.MulticastType.String() == "SQMULTICAST" {
		log.Println("STARTING TOC COMMUNICATION")
		nodeID, _ := strconv.Atoi(groupInfo.MulticastId)
		n := rand.Intn(len(connections))
		node := new(SQMulticast.NodeForSq)
		node.NodeId = uint(nodeID)
		node.Connections = connections
		node.DeliverQueue = make([]*SQMulticast.MessageSeq, 0, 100)
		node.MyConn = myConn

		SQMulticast.AppendNodes(*node)
		seq := new(SQMulticast.Sequencer)
		seq.Node = *node
		seq.SeqPort = connections[n]
		seq.Connections = connections
		seq.LocalTimestamp = 0
		seq.MessageQueue = make([]SQMulticast.MessageSeq, 0, 100)
		SQMulticast.SetSequencer(*seq)
		pool.Pool.InitThreadPool(connections, 5, util.SQMULTICAST, nil, *port)
		log.Println("The sequencer nodes is at port", seq.SeqPort.Connection.Target())
	}

	if groupInfo.MulticastType.String() == "SCMULTICAST" {
		log.Println("STARTING TOD COMMUNICATION")
		pool.Pool.InitThreadPool(connections, 5, util.SCMULTICAST, nil, *port)
		go impl.Receive(connections, *port)
	}
	if groupInfo.MulticastType.String() == "VCMULTICAST" {
		log.Println("STARTING CO COMMUNICATION")
		var wg sync.Mutex
		var myNode int32

		wg.Lock()
		VectorClockMulticast.InitLocalTimestamp(&wg, len(connections))
		VectorClockMulticast.SetMyNode(myNode)
		VectorClockMulticast.SetConnections(connections)
		pool.Pool.InitThreadPool(connections, 5, util.VCMULTICAST, nil, *port)

	}
	return nil

}

//Change status of the member
func ChangeStatus(groupInfo *ServiceProto.Group, multicastGroup *MulticastGroup, status ServiceProto.Status) (*ServiceProto.Group, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = RegClient.GetStatus(context.Background(), &ServiceProto.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}
		update(groupInfo, multicastGroup)
	}

	if groupInfo.Status == ServiceProto.Status_CLOSED || groupInfo.Status == ServiceProto.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}

	return groupInfo, nil
}

//update status of member
func update(groupInfo *ServiceProto.Group, multicastGroup *MulticastGroup) {
	multicastGroup.groupMu.Lock()
	defer multicastGroup.groupMu.Unlock()

	multicastGroup.Group.Status = ServiceProto.Status_name[int32(groupInfo.Status)]

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
