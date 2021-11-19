package restApi

import (
	"errors"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/ServiceProto"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/pool"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"sync"
)

type routes struct {
	router *gin.Engine
}

func getMessages(c *gin.Context) {

	mId := c.Param("mId")
	group, ok := MulticastGroups[mId]
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "group" + mId + "not found"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"messages": group.Messages})

}

func getInfoGroup(c *gin.Context) {

	groups := make([]*MulticastInfo, 0)

	for _, group := range MulticastGroups {
		group.groupMu.RLock()
		groups = append(groups, group.Group)
		group.groupMu.RUnlock()
	}

	c.IndentedJSON(http.StatusOK, gin.H{"groups": groups})

}

func addGroup(c *gin.Context) {

	var req MulticastReq
	err := c.BindJSON(&req)

	multicastId := req.MulticastId

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	multicastType, ok := ServiceProto.TypeMulticast_value[req.MulticastType]
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "multicast type not supported"})
	}
	log.Println("Creating multicast with type ", ServiceProto.TypeMulticast(multicastType))

	GMu.Lock()
	defer GMu.Unlock()

	group, _ := MulticastGroups[multicastId]

	registrationAns, err := RegClient.Register(context.Background(), &ServiceProto.RegInfo{
		MulticastId:   multicastId,
		MulticastType: ServiceProto.TypeMulticast(multicastType),
		Port:          uint32(GrpcPort),
	})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	members := make(map[string]Member, 0)

	for memberId, member := range registrationAns.GroupInfo.Members {
		members[memberId] = Member{
			MemberId: member.Id,
			Address:  member.Address,
			Ready:    member.Ready,
		}
	}

	group = &MulticastGroup{
		clientId: registrationAns.ClientId,
		Group: &MulticastInfo{
			MulticastId:      registrationAns.GroupInfo.MulticastId,
			MulticastType:    ServiceProto.TypeMulticast(multicastType).String(),
			ReceivedMessages: 0,
			Status:           ServiceProto.Status_name[int32(registrationAns.GroupInfo.Status)],
			Members:          members,
		},
		Messages: make([]Message, 0),
		groupMu:  sync.RWMutex{},
	}

	MulticastGroups[registrationAns.GroupInfo.MulticastId] = group

	go InitGroup(registrationAns.GroupInfo, group, GrpcPort)

	c.IndentedJSON(http.StatusOK, gin.H{"group": group.Group})

}

func sendMessage(c *gin.Context) {

	mId := c.Param("mId")
	var req Message
	err := c.BindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error input"})
	}

	group, ok := MulticastGroups[mId]
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "group" + mId + "not found"})
	}

	group.groupMu.RLock()
	defer group.groupMu.RUnlock()

	if ServiceProto.Status(ServiceProto.Status_value[group.Group.Status]) != ServiceProto.Status_ACTIVE {
		c.IndentedJSON(http.StatusTooEarly, gin.H{"error": "group not ready"})
		return
	}
	log.Println("Trying to multicasting message to group ", mId)
	//multicastType := group.Group.MulticastType
	payload := req.Payload
	//mtype, ok := ServiceProto.TypeMulticast_value[multicastType]

	log.Println("Trying to sending ", payload)

	pool.Pool.Message <- &rpc.Packet{Message: payload}

	c.IndentedJSON(http.StatusOK, gin.H{"message": payload})
}

func closeGroup(c *gin.Context) {

	mId := c.Param("mId")
	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[mId]

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "group not found"})
	}

	_, err := ServiceProto.RegistryClient.CloseGroup(context.Background(), &ServiceProto.RequestData{
		MulticastId: group.Group.MulticastId,
		ClientId:    mId,
	}, nil)
	if err != nil {
		log.Println("Error in closing group")
	}

}

func startGroup(c *gin.Context) {

	mId := c.Param("mId")

	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[mId]

	if !ok {
		statusCode := http.StatusInternalServerError
		c.IndentedJSON(statusCode, gin.H{"error": errors.New("groups doesn't exist")})
	}

	info, err := RegClient.StartGroup(context.Background(), &ServiceProto.RequestData{
		MulticastId: group.Group.MulticastId,
		ClientId:    group.clientId})

	if err != nil {
		log.Println("Error in start group ", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": info})

}
