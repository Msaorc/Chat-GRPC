package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Msaorc/Chat-GRPC/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to the server: %v", err)
	}

	defer connection.Close()

	chatClient := proto.NewChatClient(connection)

	stream, err := chatClient.Join(context.Background())
	if err != nil {
		log.Fatalf("Could not join the server: %v", err)
	}

	fmt.Print("Enter your name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	user := scanner.Text()

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error receiving message: %v", err)
			}

			fmt.Printf("[%s] %s: %s\n", time.Unix(msg.Timestamp, 0).Format("15:04:05"), msg.User, msg.Message)
		}
	}()

	for scanner.Scan() {
		msg := &proto.Message{
			User:      user,
			Message:   scanner.Text(),
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(msg); err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
	}
}
