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
	log.Println("Create New Accont: ", req.Name)
	response := &pb.CreateAccountResponse{
		ID: "abc1",
	}

	return response, nil
}

func (bs *BankingServer) ReadAccount(ctx context.Context, req *pb.ReadAccountRequest) (*pb.ReadAccountResponse, error) {
	log.Println("Read Data:", req.ID)
	response := &pb.ReadAccountResponse{}
	return response, nil
}

func (bs *BankingServer) DepositMoney(ctx context.Context, req *pb.DepositMoneyRequest) (*pb.DepositMoneyResponse, error) {
	log.Println("Account ID: ", req.ID, " Deposit Money: ", req.Money)
	response := &pb.DepositMoneyResponse{
		Result: "success",
	}
	return response, nil
}

func (bs *BankingServer) WithdrawMoney(ctx context.Context, req *pb.WithdrawMoneyRequest) (*pb.WithdrawMoneyResponse, error) {
	log.Println("Account ID: ", req.ID, " Withdraw Money: ", req.Money)
	response := &pb.WithdrawMoneyResponse{
		Result: "success",
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
