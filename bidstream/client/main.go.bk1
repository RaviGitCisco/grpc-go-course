package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"
	"google.golang.org/grpc"
)

//type client struct{}

func startStreaming(client pb.BidirectionalService_BidirectionalStreamClient, interval int) {
	for {
		req := &pb.Request{Content: "Client request"}
		if err := client.Send(req); err != nil {
			log.Printf("Error sending request: %v", err)
			return
		}

		for {
			res, err := client.Recv()
			if err == io.EOF {
				log.Println("Server closed the connection.")
				return
			}
			if err != nil {
				log.Printf("Error receiving response: %v", err)
				return
			}
			fmt.Printf("Received response: %s\n", res.Content)
			time.Sleep(5 * time.Second) // Adjust interval as needed
		}

	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <interval_seconds>")
		os.Exit(1)
	}

	interval, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid interval: %v", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBidirectionalServiceClient(conn)
	stream, err := client.BidirectionalStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	startStreaming(stream, interval)
}
