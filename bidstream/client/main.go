package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <interval_seconds> ")
		os.Exit(1)
	}

	interval, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid interval: %v", err)
	}

	dieafter, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid message count: %v", err)
	}

	log.Printf("interval %v", interval)
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBidirectionalServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), (time.Duration((uint32)(dieafter)) * time.Second))
	defer cancel()

	stream, err := client.BidirectionalStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}
	defer stream.CloseSend()

	log.Printf("Stream created")
	// Set up a goroutine to send RecordRequest messages every 60 seconds
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("connection closed.")
				return
			case <-time.After(time.Duration((uint32)(interval)) * time.Second):
				err = stream.Send(&pb.Request{
					Content: "Clients request",
				})
				log.Printf("Sent message in stream.")
			}
		}
	}()

	// Receive and print RecordResponse messages
	for {
		response, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive response: %v", err)
		}
		fmt.Printf("Received RecordResponse: %s\n", response.String())
	}
}
