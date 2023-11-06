package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"

	pb "github.com/dang252/Golang-gRPC-Banking/Bankingpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	Id int32
	Name        string
	Email       string
	PhoneNumber string
}

type UserManagement struct {
	pb.UnimplementedUserManagementServer
}

func init() {
	DatabaseConnection()
}

var DB *gorm.DB
var err error

func DatabaseConnection() {
	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("db connection error: ", err)
	}
	log.Println("db connection successful")
}

func (ums *UserManagement) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("Create New Accont: ", req.Name)
	user := User{Name: req.Name, Email: req.Email, PhoneNumber: req.PhoneNumber}
	res :=DB.Create(&user)
	if res.RowsAffected == 0 {
   		return nil, errors.New("error create user")
 	}
	response := &pb.CreateUserResponse{
		Id: user.Id,
	}

	return response, nil
}

func (ums *UserManagement) ReadUser(ctx context.Context, req *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	user := User{Id: req.Id}
	log.Println("Read User:", req.Id)
	res:=DB.First(&user);
	if res.RowsAffected == 0 {
   		return nil, status.Error(codes.InvalidArgument, "invalid Id")
 	}
	response := &pb.ReadUserResponse{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
	return response, nil
}

func (ums *UserManagement) UserReport(ctx context.Context, req *pb.UserReportRequest) (*pb.UserReportResponse, error) {
	banking_conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Proplem with banking server: %v", err)
	}
	defer banking_conn.Close()

	if err != nil {
		return nil, errors.New("Can't connect to banking server")
	}

	banking_client := pb.NewBankingServiceClient(banking_conn)

	accounts, err := banking_client.GetUserAllAccount(ctx, &pb.GetUserAllAccountRequest{Id : req.Id})

	if err != nil {
		return nil, errors.New("Can't get account data")
	} else if len(accounts.BankAccountIds) == 0 {
		return &pb.UserReportResponse{Message: "You have no account!"}, nil
	} 

	reports := []*pb.BankAccountReportResponse{}
	for _, id := range accounts.BankAccountIds {
		report, err := banking_client.BankAccountReport(ctx, &pb.BankAccountReportRequest{AccountId: id}); 
		if err != nil {
			return nil, errors.New("Can't get report")
		}
		reports = append(reports, report)
	}
	response := &pb.UserReportResponse{
		BankAccountReport: reports,
	}
	return response, nil
}

func (ums *UserManagement) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	banking_conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Proplem with banking server: %v", err)
	}
	defer banking_conn.Close()

	if err != nil {
		return nil, errors.New("Can't connect to banking server")
	}

	banking_client := pb.NewBankingServiceClient(banking_conn)

	accounts, err := banking_client.GetUserAllAccount(ctx, &pb.GetUserAllAccountRequest{Id : req.Id})

	if err != nil {
		return nil, errors.New("Can't get account data")
	} 

	for _, id := range accounts.BankAccountIds {
		banking_client.DeleteBankAccount(ctx, &pb.DeleteBankAccountRequest{Id: id}); 
	}
	user := &User{}
	if err := DB.Where(&User{ Id : req.Id}).Delete(&user).Error; err!=nil{
		return nil, errors.New("Can't get user accounts")
	}
	response := &pb.DeleteUserResponse{
		Result : "Account deleted",
	}
	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	log.Printf("listening at: %v", lis.Addr())

	s := grpc.NewServer()

	// pb.RegisterBankingServiceServer(s, &BankingServer{})
	pb.RegisterUserManagementServer(s, &UserManagement{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
