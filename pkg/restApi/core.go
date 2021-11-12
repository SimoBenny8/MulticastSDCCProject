package restApi

import (
	"MulticastSDCCProject/pkg/ServiceRegistry/ServiceProto"
	"MulticastSDCCProject/pkg/ServiceRegistry/client"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
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

func Run(grpcP, restPort uint, registryAddr, relativePath string, numThreads, dl uint, debug bool) error {
	GrpcPort = grpcP
	Delay = dl

	var err error
	Application = true
	//utils2.Vectorclock = utils2.NewVectorClock(4)

	RegClient, err = client.Connect(registryAddr)

	if err != nil {
		return err
	}

	//utils.GoPool.Initialize(int(numThreads))

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
	groups.GET("/deliver/:mId", getDeliverQueue)
}
