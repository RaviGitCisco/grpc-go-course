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
	pb.BidirectionalServiceServer
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
				log.Println("Close channel message received.")
				return
			default:
				msg, err := stream.Recv()
				if err == io.EOF {
					log.Println("Stream received EOF message.")
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
				log.Println("Closing the Send channel.")
				return
			case msg := <-s.sendCh:
				if msg == nil {
					// Log an error or handle the nil message case appropriately
					log.Printf("Error: Attempted to send nil message.")
					return
				}
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
		log.Println("Timer started to receive messages from client")
		for {
			select {
			case <-timer.C:
				log.Println("Connection timed out. Closing.")
				s.closeConnection()
				return
			case <-s.recvCh:
				// Reset the timer if a message is received
				if !timer.Stop() {
					log.Println("Stop timer block.")
					<-timer.C
				}
				log.Println("Resetting the timer")
				timer.Reset(s.timeout)
			case <-s.closeCh:
				log.Println("Closing the channel")
				return
			}
			log.Println("I am looping here.")
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	return nil
}

func (s *server) closeConnection() {
	s.closeOnce.Do(func() {
		log.Println("Closing the channel")
		close(s.closeCh)
	})
}

// The following methods are required to implement BidirectionalServiceServer interface
//func (s *server) BidirectionalServiceServer(ctx context.Context, request *pb.Message) (*pb.Message, error) {
// Implement the method logic here
//return &pb.Message{}, nil
//}

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
		log.Println("Listner routine started.")
	}()

	// Simulate sending messages to the server
	go func() {
		for i := 0; i < 5; i++ {
			msg := &pb.Message{Content: fmt.Sprintf("Server response %d", i)}
			serverInstance.sendCh <- msg
			log.Printf("Sending message %v", msg.Content)
		}
		close(serverInstance.sendCh)
		log.Println("Exiting send routine.")
	}()

	// Wait for user input to exit
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
	close(serverInstance.sendCh)
	close(serverInstance.recvCh)
	srv.GracefulStop()
}
