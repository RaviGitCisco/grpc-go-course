package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.BidirectionalServiceServer
	sendCh  chan *pb.Message
	timeout time.Duration
	closeCh chan struct{}
}

func (s *server) BidirectionalStream(stream pb.BidirectionalService_BidirectionalStreamServer) error {
	s.sendCh = make(chan *pb.Message)
	s.closeCh = make(chan struct{})
	var wg sync.WaitGroup
	defer close(s.sendCh)
	defer close(s.closeCh)

	timer := time.NewTimer(s.timeout)
	defer timer.Stop()
	log.Printf("Timeout value %v", s.timeout)

	// Handle incoming messages
	log.Printf("In BidirectionalStream starting the receive goroutine.")
	wg.Add(1)
	go func() {
		log.Printf("Receive goroutine started.")
		defer wg.Done()

		for {
			log.Println("Looping in receive routine for loop")
			select {
			case <-s.closeCh:
				log.Printf("Close channel message received. Exiting receive routine.")
				return
			default:
				msg, err := stream.Recv()
				if err == io.EOF {
					log.Printf("Stream received EOF message.")
					close(s.closeCh)
					log.Println("Exiting receive routine.")
					return
				}
				if err != nil {
					log.Printf("Error receiving message: %v", err)
					close(s.closeCh)
					log.Println("Closed closeCh, exiting receive routine.")
					return
				}
				log.Printf("Received message: %v\n", msg)
				if !timer.Stop() {
					<-timer.C
				}
				log.Println("Timer running receive routine")
				timer.Reset(s.timeout)
				log.Println("Timer reset in receive routine")
			}
			log.Println("How its in receive routine.")
		}
	}()

	// Handle outgoing messages
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-s.closeCh:
				log.Printf("Exiting the Send routine.")
				return
			case msg := <-s.sendCh:
				if err := stream.Send(msg); err != nil {
					log.Printf("Error sending message: %v", err)
					close(s.closeCh)
					return
				}
				log.Printf("Message sent successfully to Client.")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var i int
		for {
			select {
			case <-s.closeCh:
				log.Println("Closing the Message generator routine")
				return
			default:
				i += 1
				msg := &pb.Message{Content: fmt.Sprintf("Server response %d\n", i)}
				log.Printf("Sending message %v", msg.Content)
				s.sendCh <- msg
				log.Println("Message sent to sendch.")
				time.Sleep(10 * time.Second)
			}
		}
	}()

	// Monitor timeout for receiving messages
	wg.Add(1)
	go func() {
		log.Println("Timer started to receive messages from client")
		for {
			select {
			case <-timer.C:
				log.Println("Connection timed out. Closing closeCh.")
				close(s.closeCh)
				log.Println("Exiting timeout routine.")
				return
			case <-s.closeCh:
				log.Println("Closing the timer channel")
				return
			}
		}
	}()

	wg.Wait()
	log.Printf("All goroutines returned")
	return nil
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

	serverInstance := &server{
		timeout: 10 * time.Second,
	}

	pb.RegisterBidirectionalServiceServer(srv, serverInstance)

	// Continue prompting the user until they type 'exit'
	for {
		go func() {
			log.Println("Server is listening on :50051")
			if err := srv.Serve(listener); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
			log.Println("Listner routine started.")
		}()
		// Wait for user input to exit
		fmt.Println("Type 'exit' to quit the program.")

		// Create a scanner to read user input from the console
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Enter text: ")
		scanner.Scan()
		input := scanner.Text()

		if strings.ToLower(input) == "exit" {
			fmt.Println("Exiting the program.")
			break
		}
	}
	log.Println("About to exit.")
	//close(serverInstance.sendCh)
	//srv.GracefulStop()
	log.Println("Exiting main")
}
