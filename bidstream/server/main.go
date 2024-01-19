package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"
	"google.golang.org/grpc"
)

type server struct {
	sendCh    chan *pb.Message
	recvCh    chan *pb.Message
	timeout   time.Duration
	closeOnce sync.Once
	closeCh   chan struct{}
}

func (s *server) BidirectionalStream(stream pb.BidirectionalService_BidirectionalStreamServer) error {
	var wg sync.WaitGroup

	// Handle incoming messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-s.closeCh:
				return
			default:
				msg, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("Error receiving message: %v", err)
					s.closeConnection()
					return
				}
				fmt.Printf("Received message: %s\n", msg.Content)
				s.recvCh <- msg
			}
		}
	}()

	// Handle outgoing messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-s.closeCh:
				return
			case msg := <-s.sendCh:
				if err := stream.Send(msg); err != nil {
					log.Printf("Error sending message: %v", err)
					s.closeConnection()
					return
				}
			}
		}
	}()

	// Monitor timeout for receiving messages
	go func() {
		timer := time.NewTimer(s.timeout)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				log.Println("Connection timed out. Closing.")
				s.closeConnection()
				return
			case <-s.recvCh:
				// Reset the timer if a message is received
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(s.timeout)
			case <-s.closeCh:
				return
			}
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	return nil
}

func (s *server) closeConnection() {
	s.closeOnce.Do(func() {
		close(s.closeCh)
	})
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	sendCh := make(chan *pb.Message)
	recvCh := make(chan *pb.Message)
	closeCh := make(chan struct{})

	serverInstance := &server{
		sendCh:  sendCh,
		recvCh:  recvCh,
		timeout: 120 * time.Second,
		closeCh: closeCh,
	}

	pb.RegisterBidirectionalServiceServer(srv, serverInstance)

	go func() {
		log.Println("Server is listening on :50051")
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Simulate sending messages to the server
	go func() {
		for i := 0; i < 5; i++ {
			serverInstance.sendCh <- &pb.Message{Content: fmt.Sprintf("Server response %d", i)}
		}
		close(serverInstance.sendCh)
	}()

	// Wait for user input to exit
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
	close(listener)
	close(serverInstance.sendCh)
	close(serverInstance.recvCh)
}
