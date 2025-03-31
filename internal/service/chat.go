package service

import (
	"fmt"
	"sync"

	"github.com/Msaorc/Chat-GRPC/internal/proto"
)

type ChatService struct {
	proto.UnimplementedChatServer
	mu       sync.Mutex
	clients  map[proto.Chat_JoinServer]bool
	messages chan *proto.Message
}

func NewChatService() *ChatService {
	return &ChatService{
		clients:  make(map[proto.Chat_JoinServer]bool),
		messages: make(chan *proto.Message),
	}
}

func (c ChatService) Join(stream proto.Chat_JoinServer) error {
	c.mu.Lock()
	c.clients[stream] = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.clients, stream)
		c.mu.Unlock()
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}

			c.messages <- msg
		}
	}()

	for messages := range c.messages {
		c.mu.Lock()
		for client := range c.clients {
			if err := client.Send(messages); err != nil {
				fmt.Printf("Error trying to send message %v", err)
			}
		}
		c.mu.Unlock()
	}

	return nil
}
