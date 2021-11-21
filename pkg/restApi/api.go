package restApi

import (
	"context"
	"errors"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/proto"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/pool"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/rpc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
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
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var req MulticastReq
	err := c.BindJSON(&req)

	multicastId := req.MulticastId

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	multicastType, ok := proto.TypeMulticast_value[req.MulticastType]
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "multicast type not supported"})
	}
	log.Println("Creating multicast with type ", proto.TypeMulticast(multicastType))

	GMu.Lock()
	defer GMu.Unlock()

	group, _ := MulticastGroups[multicastId]

	registrationAns, err := RegClient.Register(ctx, &proto.RegInfo{
		MulticastId:   multicastId,
		MulticastType: proto.TypeMulticast(multicastType),
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
			MulticastType:    proto.TypeMulticast(multicastType).String(),
			ReceivedMessages: 0,
			Status:           proto.Status_name[int32(registrationAns.GroupInfo.Status)],
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

	if proto.Status(proto.Status_value[group.Group.Status]) != proto.Status_ACTIVE {
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
	context_, _ := context.WithTimeout(c, time.Second)
	mId := c.Param("mId")
	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[mId]

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "group not found"})
	}

	_, err := RegClient.CloseGroup(context_, &proto.RequestData{
		MulticastId: group.Group.MulticastId,
		ClientId:    group.clientId,
	})
	log.Println("Group Closed")
	if err != nil {
		log.Println("Error in closing group", err)
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

	info, err := RegClient.StartGroup(context.Background(), &proto.RequestData{
		MulticastId: group.Group.MulticastId,
		ClientId:    group.clientId})

	if err != nil {
		log.Println("Error in start group ", err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": info})

}
