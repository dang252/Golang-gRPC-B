package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	pb "github.com/dang252/Golang-gRPC-Banking/Bankingpb"
)

var (
	// command-line options:
	// gRPC server endpoint
	BankingServerEndpoint = flag.String("Banking-server-endpoint", "localhost:50051", "Banking server endpoint")
	UserManagementServerEndpoint = flag.String("User-management-server-endpoint", "localhost:50052", "User management server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	log.Println("proxy running")
	err := pb.RegisterBankingServiceHandlerFromEndpoint(ctx, mux, *BankingServerEndpoint, opts)
	err1 := pb.RegisterUserManagementHandlerFromEndpoint(ctx, mux, *UserManagementServerEndpoint, opts)
	if err != nil {
		return err
	}
	if err1 != nil {
		return err
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
