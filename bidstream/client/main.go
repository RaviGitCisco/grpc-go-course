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
	sendCh   chan *pb.Message
	recvCh   chan *pb.Message
	closeCh  chan struct{}
	sendDone chan struct{}
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
		sendCh:   make(chan *pb.Message),
		recvCh:   make(chan *pb.Message),
		closeCh:  make(chan struct{}),
		sendDone: make(chan struct{}),
	}

	client := pb.NewBidirectionalServiceClient(conn)

	stream, err := client.BidirectionalStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	var sendWg sync.WaitGroup
	var recvWg sync.WaitGroup

	// Handle outgoing messages
	sendWg.Add(1)
	go func() {
		log.Printf("In send goroutine")
		defer sendWg.Done()
		defer close(clientInstance.sendCh)
		for i := -1; i < messageCount; i++ {
			msg := &pb.Message{Content: fmt.Sprintf("Client message %d", i)}
			log.Printf("Send Message %v", msg)
			log.Printf("Calling stream.send")
			if err := stream.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
			log.Println("Message sent successfully to Client.")
			log.Printf("Sent message sleeping %v seconds", (time.Duration(interval) * time.Second))
			select {
			case <-time.After(time.Duration(interval) * time.Second):
			case <-clientInstance.closeCh:
				return
			}
		}
	}()

	// Wait for the send goroutine to complete
	go func() {
		sendWg.Wait()
		close(clientInstance.closeCh)
	}()

	// Handle incoming messages
	recvWg.Add(1)
	go func() {
		defer recvWg.Done()
		defer close(clientInstance.recvCh)
		for {
			log.Println("Waiting for response from server.")
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Println("Receving gorountine reached EOF.")
				return
			}
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			fmt.Printf("Received message: %s\n", msg.Content)
			clientInstance.recvCh <- msg
		}
	}()

	// Additional goroutine for processing received messages
	go func() {
		for receivedMsg := range clientInstance.recvCh {
			// Process receivedMsg
			fmt.Printf("Processing received message: %s\n", receivedMsg.Content)
		}
	}()

	// Wait for both goroutines to finish
	sendWg.Wait()

	// Close the stream properly after sending all messages
	if err := stream.CloseSend(); err != nil {
		log.Printf("Error closing the stream: %v", err)
	}

	// Wait for the server responses
	recvWg.Wait()

	// Signal that the sending is done
	close(clientInstance.sendDone)
}
