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

func (s *server) BidirectionalStream(stream pb.BidirectionalService_BidirectionalStreamServer) error {
	// Start a goroutine to handle incoming events from other modules
	/*
		go func() {
			for {
				select {
				case event := <-s.eventChannel:
					// Send RecordResponse to the client when an event occurs
					if err := stream.Send(&pb.Response{Content: event}); err != nil {
						log.Printf("Error sending RecordResponse: %v", err)
					}
				}
			}
		}()
	*/
	// Start a goroutine to simulate events from other modules
	log.Printf("Created timer routine")
	go func() {
		for {
			time.Sleep(2 * time.Second) // Simulating an event occurring every 5 seconds
			select {
			case value, ok := <-s.eventChannel:
				if !ok {
					log.Printf("Channel is closed %v", value)
					return
				}
			default:
				s.eventChannel <- "Event occurred"
			}
		}
	}()
	log.Printf("Inside Server")
	go func() {
		for event := range s.eventChannel {
			if event == "Exit" {
				log.Printf("Exiting Server. Client liveness expired.")
				close(s.eventChannel)
				return
			}
			log.Printf("Sending message in stream")
			// Send RecordResponse to the client when an event occurs
			if err := stream.Send(&pb.Response{Content: event}); err != nil {
				log.Printf("Error sending RecordResponse: %v", err)
			}
		}
	}()
	// Receive RecordRequest from the client
	for {
		log.Printf("Waiting for request from client")
		_, err := stream.Recv()
		if err != nil {
			// Handle client disconnection or errors
			return err
		}

		value, ok := <-s.eventChannel

		if !ok {
			log.Printf("Channel is closed %v.", value)
			return err
		}
		// Update the last received timestamp
		s.mu.Lock()
		log.Printf("Last received %v.", time.Now())
		s.lastReceived = time.Now()
		s.mu.Unlock()
	}
}

func (s *server) checkLiveliness() {
	// Check liveliness every 60 seconds
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		elapsed := time.Since(s.lastReceived)
		log.Printf("Last received %v elapsed %v", s.lastReceived, elapsed)
		s.mu.Unlock()

		if elapsed > 40*time.Second {
			log.Println("Client is not sending RecordRequest for more than 120 seconds. Resetting connection.")
			// You can perform any cleanup or reset connection logic here
			s.eventChannel <- "Exit"
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	acctserver := grpc.NewServer()
	recordService := newRecordServer()
	pb.RegisterBidirectionalServiceServer(acctserver, recordService)

	server := &server{}
	// Start the server
	wg.Add(2)
	log.Printf("Serving now.")
	go func() {
		defer wg.Done()
		log.Printf("Serving now.")
		acctserver.Serve(listener)
		log.Printf("Service Server returned.")
	}()

	// Start a goroutine to check client liveliness
	go func() {
		defer wg.Done()
		recordService.checkLiveliness()
		log.Printf("CheckLiveliness returned")
		server.eventChannel <- "Exit"
	}()

	log.Println("Server started at :50051")
	wg.Wait()
}
