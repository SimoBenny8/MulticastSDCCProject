package server

import (
	"context"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"strings"
	"sync"
)

type RegistryServer struct {
	proto.UnimplementedRegistryServer
}

var (
	groups   map[string]*MulticastGroup
	Mugroups sync.RWMutex
)

type MulticastGroup struct {
	mu        sync.RWMutex
	groupInfo *proto.Group
}

func init() {
	groups = make(map[string]*MulticastGroup)
}

// Registration of the node in a specific group
func (s *RegistryServer) Register(ctx context.Context, in *proto.RegInfo) (*proto.RegAnswer, error) {

	log.Println("Starting register service ")

	source, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}
	src := source.Addr.String()
	srcAddr := src[:strings.LastIndexByte(src, ':')]

	srcAddr = fmt.Sprintf("%s:%d", srcAddr, in.Port)
	//utils.MyAdress = srcAddr
	log.Println("Registration of the group ", in.MulticastId, "with client", srcAddr)
	multicastId := in.MulticastId
	Mugroups.Lock()
	defer Mugroups.Unlock()
	// Checking if the group already exists
	multicastGroup, exists := groups[multicastId]

	if exists {
		log.Println("group already exists")
		return nil, nil
	} else {
		// Creating the group
		multicastGroup = &MulticastGroup{groupInfo: &proto.Group{
			MulticastId:   multicastId,
			MulticastType: in.MulticastType,
			Status:        proto.Status_OPENING,
			Members:       make(map[string]*proto.MemberInfo),
		}}
		groups[multicastId] = multicastGroup
	}
	// Registering node to the group
	memberInfo := new(proto.MemberInfo)
	memberInfo.Id = srcAddr
	memberInfo.Address = srcAddr
	memberInfo.Ready = false
	//Adding the MemberInfo to the multicastGroup Map
	multicastGroup.groupInfo.Members[srcAddr] = memberInfo

	return &proto.RegAnswer{
		ClientId:  srcAddr,
		GroupInfo: multicastGroup.groupInfo,
	}, nil
}

// Start enables multicast in the group.
func (s *RegistryServer) StartGroup(ctx context.Context, in *proto.RequestData) (*proto.Group, error) {
	_, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}
	multicastId := in.MulticastId
	clientId := in.ClientId
	Mugroups.RLock()
	defer Mugroups.RUnlock()
	// Checking if the group exists
	mGroup, ok := groups[multicastId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}
	mGroup.mu.Lock()
	defer mGroup.mu.Unlock()

	if _, ok = mGroup.groupInfo.Members[clientId]; !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	mGroup.groupInfo.Status = proto.Status_STARTING
	return mGroup.groupInfo, nil
}

//Method that changes status of a member. The member is ready for communication.
func (s *RegistryServer) Ready(ctx context.Context, in *proto.RequestData) (*proto.Group, error) {
	var member *proto.MemberInfo
	_, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}
	multicastId := in.MulticastId
	clientId := in.ClientId
	Mugroups.RLock()
	defer Mugroups.RUnlock()

	mGroup, ok := groups[multicastId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}

	mGroup.mu.Lock()
	defer mGroup.mu.Unlock()

	if member, ok = mGroup.groupInfo.Members[clientId]; !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	if !member.Ready {
		member.Ready = true
		mGroup.groupInfo.ReadyMembers++
	}

	if mGroup.groupInfo.ReadyMembers == uint64(len(mGroup.groupInfo.Members)) {
		mGroup.groupInfo.Status = proto.Status_ACTIVE
	}
	return mGroup.groupInfo, nil
}

// CloseGroup closes the group.
func (s *RegistryServer) CloseGroup(ctx context.Context, in *proto.RequestData) (*proto.Group, error) {

	_, ok := peer.FromContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}

	multicastId := in.MulticastId
	clientId := in.ClientId

	Mugroups.RLock()
	defer Mugroups.RUnlock()

	mGroup, ok := groups[multicastId]

	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}

	mGroup.mu.Lock()
	defer mGroup.mu.Unlock()

	if mGroup.groupInfo.Status != proto.Status_ACTIVE && mGroup.groupInfo.Status != proto.Status_CLOSING {
		return nil, status.Errorf(codes.Canceled, "Cannot closed group")
	}

	member, ok := mGroup.groupInfo.Members[clientId]

	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	mGroup.groupInfo.Status = proto.Status_CLOSING
	log.Println("start closing group")
	if member.Ready {
		member.Ready = false
		mGroup.groupInfo.ReadyMembers--
		if mGroup.groupInfo.ReadyMembers == 0 {
			mGroup.groupInfo.Status = proto.Status_CLOSED
			log.Println("Group  closed")
		}
	}

	return mGroup.groupInfo, nil
}

func Registration(s grpc.ServiceRegistrar) (err error) {
	proto.RegisterRegistryServer(s, &RegistryServer{})
	return
}

// GetStatus returns infos about the group associated to multicastId
func (s *RegistryServer) GetStatus(ctx context.Context, in *proto.MulticastId) (*proto.Group, error) {
	_, ok := peer.FromContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}

	Mugroups.RLock()
	defer Mugroups.RUnlock()

	mGroup, ok := groups[in.MulticastId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}

	mGroup.mu.RLock()
	defer mGroup.mu.RUnlock()

	return mGroup.groupInfo, nil

}
