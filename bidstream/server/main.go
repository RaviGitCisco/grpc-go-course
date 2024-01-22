package main

import (
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.BidirectionalServiceServer
	mu           sync.Mutex
	lastReceived time.Time
	eventChannel chan string
}

func newRecordServer() *server {
	return &server{
		eventChannel: make(chan string),
	}
}

func (s *server) BidirectionalStream(stream BidirectionalService_BidirectionalStreamServer) error {
	// Start a goroutine to handle incoming events from other modules
	go func() {
		for {
			select {
			case event := <-s.eventChannel:
				// Send RecordResponse to the client when an event occurs
				if err := stream.Send(&Response{Message: event}); err != nil {
					log.Printf("Error sending RecordResponse: %v", err)
				}
			}
		}
	}()

	// Receive RecordRequest from the client
	for {
		_, err := stream.Recv()
		if err != nil {
			// Handle client disconnection or errors
			return err
		}

		// Update the last received timestamp
		s.mu.Lock()
		s.lastReceived = time.Now()
		s.mu.Unlock()
	}
}

func (s *recordServer) checkLiveliness() {
	// Check liveliness every 60 seconds
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		elapsed := time.Since(s.lastReceived)
		s.mu.Unlock()

		if elapsed > 120*time.Second {
			log.Println("Client is not sending RecordRequest for more than 120 seconds. Resetting connection.")
			// You can perform any cleanup or reset connection logic here
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	recordService := newRecordServer()
	RegisterRecordServiceServer(server, recordService)

	// Start the server
	go server.Serve(listener)

	// Start a goroutine to check client liveliness
	go recordService.checkLiveliness()

	log.Println("Server started at :50051")
	select {}
}
