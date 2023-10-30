package main

import (
	"context"
	"log"
	"net"

	pb "github.com/dang252/Golang-gRPC-Banking/Bankingpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type client struct {
	ID          int
	Name        string
	Email       string
	PhoneNumber string
	Money       int
}

var clients []client

type BankingServer struct {
	pb.UnimplementedBankingServiceServer
}

func (bs *BankingServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	log.Println("Create New Accont: ", req.Name)
	newID := len(clients)
	clients = append(clients, client{
		ID:          newID,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Money:       0,
	})
	response := &pb.CreateAccountResponse{
		ID: int32(newID),
	}

	return response, nil
}

func (bs *BankingServer) ReadAccount(ctx context.Context, req *pb.ReadAccountRequest) (*pb.ReadAccountResponse, error) {
	log.Println("Read Data:", req.ID)
	if req.ID >= int32(len(clients)) {
		return nil, status.Error(codes.InvalidArgument, "Invalid ID")
	}
	response := &pb.ReadAccountResponse{
		ID:          int32(clients[req.ID].ID),
		Name:        clients[req.ID].Name,
		Email:       clients[req.ID].Email,
		PhoneNumber: clients[req.ID].PhoneNumber,
	}
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
