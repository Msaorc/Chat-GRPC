package main

import (
	"log"
	"net"

	"github.com/Msaorc/Chat-GRPC/internal/proto"
	"github.com/Msaorc/Chat-GRPC/internal/service"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	chatServer := service.NewChatService()
	proto.RegisterChatServer(server, chatServer)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
