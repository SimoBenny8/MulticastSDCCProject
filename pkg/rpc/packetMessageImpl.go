package rpc

import (
	"context"
	"log"
)

type MessageGrpcImpl struct {
}

func (m MessageGrpcImpl) mustEmbedUnimplementedPacketServiceServer() {
	panic("implement me")
}

//implementation of service message
func (m MessageGrpcImpl) SendPacket(ctx context.Context, message *Packet) (*ResponsePacket, error) {
	log.Println("Received message : ", message.Message)
	return &ResponsePacket{}, nil
}

//MessageGrpcImpl returns the pointer to the implementation.
func NewMessageGrpcImpl() *MessageGrpcImpl {
	return &MessageGrpcImpl{}
}
