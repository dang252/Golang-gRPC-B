package main

import (
	"context"
	"log"

	pb "github.com/dang252/Golang-gRPC-Banking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateAccount(c pb.BankingServiceClient, Name string, Email string, PhoneNumber string) {
	response, err := c.CreateAccount(context.Background(), &pb.CreateAccountRequest{Name: Name, Email: Email, PhoneNumber: PhoneNumber})
	if err != nil {
		log.Fatalf("Can't Create New Account")
	}
	log.Printf("ID:%v", response.GetID())
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
	// res, err := c.CreateAccount(ctx, &pb.CreateAccountRequest{Name: "nnd", Email: "nnd@gmail", PhoneNumber: ""})

	// if err != nil {
	// 	log.Fatalf("cant create")
	// }
	// log.Printf("ID:%v", res.GetID())

	CreateAccount(c, "NND", "nnd@gmail.com", "0909090909")

}
