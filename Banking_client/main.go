package main

import (
	"context"
	"log"

	pb "github.com/dang252/Golang-gRPC-Banking/Bankingpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateAccount(c pb.UserManagermentClient, Name string, Email string, PhoneNumber string) {
	response, err := c.CreateUser(context.Background(), &pb.CreateUserRequest{Name: Name, Email: Email, PhoneNumber: PhoneNumber})
	if err != nil {
		log.Fatalf("Can't Create New User")
	}
	log.Printf("ID:%v", response.GetID())
}

func ReadAccount(c pb.UserManagermentClient, ID int32) {
	response, err := c.ReadUser(context.Background(), &pb.ReadUserRequest{ID: ID})
	if err != nil {
		log.Println(err)
	}
	log.Println("User Data: ", response)
}

func DepositMoney(uc pb.BankingServiceClient, ID int32, Money int64) {
	response, err := uc.DepositMoney(context.Background(), &pb.DepositMoneyRequest{ID: ID, Money: Money})
	if err != nil {
		log.Println(err)
	}
	log.Println("result ", response)
}

func WithdrawMoney(uc pb.BankingServiceClient, ID int32, Money int64) {
	response, err := uc.WithdrawMoney(context.Background(), &pb.WithdrawMoneyRequest{ID: ID, Money: Money})
	if err != nil {
		log.Println(err)
	}
	log.Println("result ", response)
}

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Proplem with server: %v", err)
	}
	defer conn.Close()

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// defer cancel()
	c := pb.NewBankingServiceClient(conn)
	uc := pb.NewUserManagermentClient(conn)
	// res, err := c.CreateAccount(ctx, &pb.CreateAccountRequest{Name: "nnd", Email: "nnd@gmail", PhoneNumber: ""})

	// if err != nil {
	// 	log.Fatalf("cant create")
	// }
	// log.Printf("ID:%v", res.GetID())

	CreateAccount(uc, "NND", "nnd@gmail.com", "0909090909")

	CreateAccount(uc, "NND2", "nnd2@gmail.com", "0909090909")

	CreateAccount(uc, "NND3", "nnd3@gmail.com", "0909090909")

	ReadAccount(uc, 1)

	ReadAccount(uc, 0)

	DepositMoney(c, 1, 100)

	WithdrawMoney(c, 1, 50)

	WithdrawMoney(c, 1, 100)

	// ReadAccount(c, 3)

}
