package main

import (
	"context"
	"log"
	"net"

	pb "github.com/dang252/Golang-gRPC-Banking"
	"google.golang.org/grpc"
)

type BankingServer struct {
	pb.UnimplementedBankingServiceServer
}

func (bs *BankingServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	response := &pb.CreateAccountResponse{
		ID: "abc1",
	}

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	log.Printf("listening at: %v", lis.Addr())

	s := grpc.NewServer()

	pb.RegisterBankingServiceServer(s, &BankingServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
