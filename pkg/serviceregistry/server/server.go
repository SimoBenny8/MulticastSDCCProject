package server

import (
	"MulticastSDCCProject/pkg/serviceregistry/ServiceProto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"strings"
	"sync"
	"time"
)

type RegistryServer struct {
	ServiceProto.UnimplementedRegistryServer
}

var (
	groups   map[string]*Mgroup
	Mugroups sync.RWMutex
)

type Mgroup struct {
	mu        sync.RWMutex
	groupInfo *ServiceProto.Group
}

func init() {
	groups = make(map[string]*Mgroup)
}

// Registration of the node in a specific group
func (s *RegistryServer) Register(ctx context.Context, in *ServiceProto.RegInfo) (*ServiceProto.RegAnswer, error) {
	timeout := time.Second
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	log.Println("Start registering..")
	source, ok := peer.FromContext(c)
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
		if multicastGroup.groupInfo.MulticastType != in.MulticastType {
			return nil, status.Errorf(codes.InvalidArgument, "MulticastType")
		}
		if multicastGroup.groupInfo.Status != ServiceProto.Status_OPENING {
			return nil, status.Errorf(codes.PermissionDenied, "ServiceStarted")
		}
		if _, ok := multicastGroup.groupInfo.Members[srcAddr]; ok {
			return nil, status.Errorf(codes.AlreadyExists, "Already Registered")
		}
	} else {
		log.Println("Start registering ..")
		// Creating the group
		multicastGroup = &Mgroup{groupInfo: &ServiceProto.Group{
			MulticastId:   multicastId,
			MulticastType: in.MulticastType,
			Status:        ServiceProto.Status_OPENING,
			Members:       make(map[string]*ServiceProto.MemberInfo),
		}}
		groups[multicastId] = multicastGroup
	}
	// Registering node to the group
	//Creating new MemberInfo
	memberInfo := new(ServiceProto.MemberInfo)
	memberInfo.Id = srcAddr
	memberInfo.Address = srcAddr
	memberInfo.Ready = false
	//Adding the MemberInfo to the multicastGroup Map
	multicastGroup.groupInfo.Members[srcAddr] = memberInfo

	return &ServiceProto.RegAnswer{
		ClientId:  srcAddr,
		GroupInfo: multicastGroup.groupInfo,
	}, nil
}

// Start enables multicast in the group.
func (s *RegistryServer) StartGroup(ctx context.Context, in *ServiceProto.RequestData) (*ServiceProto.Group, error) {
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

	if mGroup.groupInfo.Status != ServiceProto.Status_OPENING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}

	if _, ok = mGroup.groupInfo.Members[clientId]; !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	mGroup.groupInfo.Status = ServiceProto.Status_STARTING
	return mGroup.groupInfo, nil
}

//Method that changes status of a member. The member is ready for communication.
func (s *RegistryServer) Ready(ctx context.Context, in *ServiceProto.RequestData) (*ServiceProto.Group, error) {

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

	if mGroup.groupInfo.Status != ServiceProto.Status_STARTING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}

	var member *ServiceProto.MemberInfo

	if member, ok = mGroup.groupInfo.Members[clientId]; !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	if !member.Ready {
		member.Ready = true
		mGroup.groupInfo.ReadyMembers++
	}

	if mGroup.groupInfo.ReadyMembers == uint64(len(mGroup.groupInfo.Members)) {
		mGroup.groupInfo.Status = ServiceProto.Status_ACTIVE
	}
	return mGroup.groupInfo, nil
}

// CloseGroup closes the group.
func (s *RegistryServer) CloseGroup(ctx context.Context, in *ServiceProto.RequestData) (*ServiceProto.Group, error) {

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

	if mGroup.groupInfo.Status != ServiceProto.Status_ACTIVE && mGroup.groupInfo.Status != ServiceProto.Status_CLOSING {
		return nil, status.Errorf(codes.Canceled, "Cannot closed group")
	}

	member, ok := mGroup.groupInfo.Members[clientId]

	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	mGroup.groupInfo.Status = ServiceProto.Status_CLOSING
	log.Println("start closing group")
	if member.Ready {
		member.Ready = false
		mGroup.groupInfo.ReadyMembers--
		if mGroup.groupInfo.ReadyMembers == 0 {
			mGroup.groupInfo.Status = ServiceProto.Status_CLOSED
			log.Println("Group  closed")
		}
	}

	return mGroup.groupInfo, nil
}

func Registration(s grpc.ServiceRegistrar) (err error) {
	ServiceProto.RegisterRegistryServer(s, &RegistryServer{})
	return
}

// GetStatus returns infos about the group associated to multicastId
func (s *RegistryServer) GetStatus(ctx context.Context, in *ServiceProto.MulticastId) (*ServiceProto.Group, error) {
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
