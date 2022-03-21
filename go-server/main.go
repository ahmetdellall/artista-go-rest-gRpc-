package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/ahmetdellall/grpc-examples/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost: 8080")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterTimeServiceServer(server, &srv{})
	log.Printf("starting rpc server...")
	err = server.Serve(lis)
	log.Fatal(err)

}

type srv struct {
	pb.UnimplementedTimeServiceServer
}

func (s srv) Now(ctx context.Context, request *pb.NewRequest) (*pb.TimeUpdate, error) {
	return &pb.TimeUpdate{Time: &pb.Time{
		Value: time.Now().String(),
	}}, nil
}

func (s srv) Stream(request *pb.TimeStreamRequest, server pb.TimeService_StreamServer) error {
	deadline := time.Now().Add(time.Duration(request.Length) * time.Second)
	for !time.Now().After(deadline) {
		time.Sleep(time.Millisecond * 300)
		server.Send(&pb.TimeUpdate{Time: &pb.Time{
			Value: time.Now().String(),
		}})
	}
	return nil
}
