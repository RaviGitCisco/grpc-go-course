package main

import (
	"testing"
	"time"

	pb "github.com/RaviGitCisco/grpc-go-course/bidstream/proto"
)

func Test_server_BidirectionalStream(t *testing.T) {
	type fields struct {
		BidirectionalServiceServer pb.BidirectionalServiceServer
		sendCh                     chan *pb.Message
		timeout                    time.Duration
		closeCh                    chan struct{}
	}
	type args struct {
		stream pb.BidirectionalService_BidirectionalStreamServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				BidirectionalServiceServer: tt.fields.BidirectionalServiceServer,
				sendCh:                     tt.fields.sendCh,
				timeout:                    tt.fields.timeout,
				closeCh:                    tt.fields.closeCh,
			}
			if err := s.BidirectionalStream(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.BidirectionalStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
