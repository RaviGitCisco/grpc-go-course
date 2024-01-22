package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.BidirectionalServiceServer
}

func (s *server) BidirectionalStream(stream pb.BidirectionalService_BidirectionalStreamServer) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(10 * time.Second)

	// Handle incoming requests
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				log.Println("Connection timed out. Closing.")
				cancel()
				return
			default:
				req, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("Error receiving request: %v", err)
					return
				}
				fmt.Printf("Received request: %s\n", req.Content)

				// Reset the timer since a request was received
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(10 * time.Second)

				s.sendResponses(stream)
			}
		}
	}()

	// Monitor timeout for incoming requests
	go func() {
		<-time.After(10 * time.Second)
		log.Println("Connection timed out. Closing.")
		cancel()
	}()

	// Wait for the incoming requests goroutine to finish
	wg.Wait()
	return nil
}

func (s *server) sendResponses(stream pb.BidirectionalService_BidirectionalStreamServer) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a response (or type 'exit' to close): ")
		scanner.Scan()
		response := scanner.Text()

		if response == "exit" {
			fmt.Println("Closing connection as requested.")
			return
		}

		res := &pb.Response{Content: response}
		if err := stream.Send(res); err != nil {
			log.Printf("Error sending response: %v", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	servinst := &server{}
	pb.RegisterBidirectionalServiceServer(srv, servinst)

	log.Println("Server is listening on :50051")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
