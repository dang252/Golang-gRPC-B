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

// type client struct {
// 	Id          int
// 	Name        string
// 	Email       string
// 	PhoneNumber string
// 	Money       int64
// }

// var clients []client

// type user struct {
// 	Id int32
// 	name        string
// 	email       string
// 	phone_number string
// }

type User struct {
	Id int32
	Name        string
	Email       string
	PhoneNumber string
}

// type BankAccount struct {
// 	Id int32
// 	UserId int32
// 	OpeningDate time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
//     CurrentBalance int64
// }

// type Transaction struct {
// 	Id int32
// 	AccountId int32
// 	Type string
// 	Amount int64 
// 	Date time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
// }

// type bank_account struct {
// 	Id int32
// 	user_Id int32
// 	opening_date time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
//     current_balance int64
// }

// type transaction struct {
// 	Id int32
// 	account_Id int32
// 	transaction_type string
// 	ammount int64 
// 	date time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
// }

type UserManagement struct {
	pb.UnimplementedUserManagementServer
}

// type BankingServer struct {
// 	pb.UnimplementedBankingServiceServer
// }

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

// func (bs *BankingServer) CreateBankAccount(ctx context.Context, req *pb.CreateBankAccountRequest) (*pb.CreateBankAccountResponse, error) {
// 	log.Println("Create New Bank Accont for user: ", req.UserId, " ", req.Balance)
// 	new_account := BankAccount{UserId: req.UserId, OpeningDate: time.Now(), CurrentBalance : req.Balance }
// 	res:=DB.Create(&new_account)
// 	if res.RowsAffected == 0 {
//    		return nil, errors.New("error create Bank account")
//  	}
// 	response := &pb.CreateBankAccountResponse{
// 		Id: new_account.Id,
// 	}
// 	return response, nil
// }

// func (bs *BankingServer) DepositMoney(ctx context.Context, req *pb.DepositMoneyRequest) (*pb.DepositMoneyResponse, error) {
// 	account := BankAccount{ Id: req.Id }
// 	res := DB.First(&account)
// 	if res.RowsAffected == 0 {
//    		return nil, errors.New("Can't Find Account")
//  	}
// 	log.Println("Account Id: ", req.Id, " Deposit Money: ", req.Money)
// 	trans_res := DB.Transaction(func(tx *gorm.DB) error {
// 		new_transaction := Transaction{AccountId : req.Id, Type : "Deposit", Amount : req.Money, Date : time.Now()}
// 		if err := DB.Create(&new_transaction).Error; err != nil {
//    			return err
//  		}
// 		account.CurrentBalance += req.Money
// 		if err := DB.Save(&account).Error; err != nil {
//    			return err
//  		}
// 		return nil
// 	})
// 	if trans_res != nil {
// 		return nil, errors.New("Can't make transaction")
// 	}
// 	log.Println("Account Id: ", req.Id, " Deposit Money: ", req.Money, " Money after deposit: ", account.CurrentBalance)
// 	response := &pb.DepositMoneyResponse{
// 		Result: "Deposit success",
// 	}
// 	return response, nil
// }

// func (bs *BankingServer) WithdrawMoney(ctx context.Context, req *pb.WithdrawMoneyRequest) (*pb.WithdrawMoneyResponse, error) {
// 	account := BankAccount{ Id: req.Id }
// 	res := DB.First(&account)
// 	if res.RowsAffected == 0 {
//    		return nil, errors.New("Can't Find Account")
//  	}
// 	if account.CurrentBalance < req.Money {
// 		return nil, errors.New("InvalId Amount")
// 	}
// 	log.Println("Account Id: ", req.Id, " Withdraw Money: ", req.Money)
// 	trans_res := DB.Transaction(func(tx *gorm.DB) error {
// 		new_transaction := Transaction{AccountId : req.Id, Type : "Withdraw", Amount : req.Money, Date : time.Now()}
// 		if err := DB.Create(&new_transaction).Error; err != nil {
//    			return err
//  		}
// 		account.CurrentBalance -= req.Money
// 		if err := DB.Save(&account).Error; err != nil {
//    			return err
//  		}
// 		return nil
// 	})
// 	if trans_res != nil {
// 		return nil, errors.New("Can't make transaction")
// 	}
// 	log.Println("Account Id: ", req.Id, " Withdraw Money: ", req.Money, " Money after withdraw: ", account.CurrentBalance)
// 	response := &pb.WithdrawMoneyResponse{
// 		Result: "Withdraw success",
// 	}
// 	return response, nil
// }

// func (bs *BankingServer) BankAccountReport(ctx context.Context, req *pb.BankAccountReportRequest) (*pb.BankAccountReportResponse, error){
// 	var transactions []Transaction
// 	if err := DB.Where(Transaction{ AccountId : req.AccountId}).Find(&transactions).Error; err!=nil{
// 		return nil, errors.New("Can't get Account report")
// 	}
// 	var total_deposit int64 = 0
// 	var total_withdraw int64 = 0
// 	trans_info := []string{}
// 	for _,transaction := range transactions {
// 		if transaction.Type == "Deposit" {
// 			total_deposit += transaction.Amount
// 		} else {
// 			total_withdraw += transaction.Amount
// 		}
// 		trans_info = append(trans_info, transaction.Type + ": " + strconv.FormatInt(transaction.Amount,10) + " Date: " + transaction.Date.String())
// 	}
// 	log.Println(total_deposit, total_withdraw)
	
// 	response := &pb.BankAccountReportResponse{
// 		AccountId: req.AccountId,
// 		Transactions: trans_info,
// 		TotalDeposit: total_deposit,
// 		TotalWithdraw: total_withdraw,
// 	}
// 	return response, nil
// }
// func (bs *BankingServer) AllAccountReport(ctx context.Context, req *pb.EmptyRequest) (*pb.AllAccountReportResponse, error) {
// 	var count int64 = 0
// 	if err := DB.Model(&BankAccount{}).Count(&count).Error; err != nil{
// 		return nil, errors.New("Can't get report")
// 	}
// 	reports := []*pb.BankAccountReportResponse{}
	
// 	for i:=1; i <= int(count); i++ {
// 		report, err := bs.BankAccountReport(ctx, &pb.BankAccountReportRequest{AccountId: int32(i)}); 
// 		if err != nil {
// 			return nil, errors.New("Can't get report")
// 		}
// 		reports = append(reports, report)
// 	}
// 	response := &pb.AllAccountReportResponse{
// 		BankAccountReport: reports,
// 	}
// 	return response, nil
// }

// func (bs *BankingServer) GetUserAllAccount(ctx context.Context, req *pb.GetUserAllAccountRequest) (*pb.GetUserAllAccountResponse, error) {
// 	accounts := []BankAccount{}
// 	if err := DB.Where(BankAccount{ UserId : req.Id}).Find(&accounts).Error; err!=nil{
// 		return nil, errors.New("Can't get user accounts")
// 	}
// 	ids := []int32{}
// 	for _,acc := range accounts {
// 		ids = append(ids, acc.Id)
// 	}
// 	response := &pb.GetUserAllAccountResponse {
// 		BankAccountIds: ids,
// 	}
// 	return response, nil
// }

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
