package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"
	"google.golang.org/grpc"
)

type client struct {
	sendCh chan *pb.Message
	recvCh chan *pb.Message
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client <interval_seconds> <message_count>")
		os.Exit(1)
	}

	interval, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid interval: %v", err)
	}

	messageCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid message count: %v", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	clientInstance := &client{
		sendCh: make(chan *pb.Message),
		recvCh: make(chan *pb.Message),
	}

	client := pb.NewBidirectionalServiceClient(conn)

	stream, err := client.BidirectionalStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	var sendWg sync.WaitGroup
	var recvWg sync.WaitGroup

	// Handle incoming messages
	recvWg.Add(1)
	go func() {
		defer recvWg.Done()
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			fmt.Printf("Received message: %s\n", msg.Content)
			clientInstance.recvCh <- msg
		}
		close(clientInstance.recvCh)
	}()

	// Handle outgoing messages
	sendWg.Add(1)
	go func() {
		defer sendWg.Done()
		for i := 0; i < messageCount; i++ {
			clientInstance.sendCh <- &pb.Message{Content: fmt.Sprintf("Client message %d", i)}
			time.Sleep(time.Duration(interval) * time.Second)
		}
		close(clientInstance.sendCh)
	}()

	// Wait for both goroutines to finish
	sendWg.Wait()
	close(stream, clientInstance.sendCh)
	recvWg.Wait()
}
