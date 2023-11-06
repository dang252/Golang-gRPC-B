package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/dang252/Golang-gRPC-Banking/Bankingpb"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type BankAccount struct {
	Id int32
	UserId int32
	OpeningDate time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
    CurrentBalance int64
}

type Transaction struct {
	Id int32
	AccountId int32
	Type string
	Amount int64 
	Date time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
}

type BankingServer struct {
	pb.UnimplementedBankingServiceServer
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

func (bs *BankingServer) CreateBankAccount(ctx context.Context, req *pb.CreateBankAccountRequest) (*pb.CreateBankAccountResponse, error) {
	log.Println("Create New Bank Accont for user: ", req.UserId, " ", req.Balance)
	new_account := BankAccount{UserId: req.UserId, OpeningDate: time.Now(), CurrentBalance : req.Balance }
	res:=DB.Create(&new_account)
	if res.RowsAffected == 0 {
   		return nil, errors.New("error create Bank account")
 	}
	response := &pb.CreateBankAccountResponse{
		Id: new_account.Id,
	}
	return response, nil
}

func (bs *BankingServer) DepositMoney(ctx context.Context, req *pb.DepositMoneyRequest) (*pb.DepositMoneyResponse, error) {
	account := BankAccount{ Id: req.Id }
	res := DB.First(&account)
	if res.RowsAffected == 0 {
   		return nil, errors.New("Can't Find Account")
 	}
	log.Println("Account Id: ", req.Id, " Deposit Money: ", req.Money)
	trans_res := DB.Transaction(func(tx *gorm.DB) error {
		new_transaction := Transaction{AccountId : req.Id, Type : "Deposit", Amount : req.Money, Date : time.Now()}
		if err := tx.Create(&new_transaction).Error; err != nil {
   			return err
 		}
		account.CurrentBalance += req.Money
		if err := tx.Save(&account).Error; err != nil {
   			return err
 		}
		return nil
	})
	if trans_res != nil {
		return nil, errors.New("Can't make transaction")
	}
	log.Println("Account Id: ", req.Id, " Deposit Money: ", req.Money, " Money after deposit: ", account.CurrentBalance)
	response := &pb.DepositMoneyResponse{
		Result: "Deposit success",
	}
	return response, nil
}

func (bs *BankingServer) WithdrawMoney(ctx context.Context, req *pb.WithdrawMoneyRequest) (*pb.WithdrawMoneyResponse, error) {
	account := BankAccount{ Id: req.Id }
	res := DB.First(&account)
	if res.RowsAffected == 0 {
   		return nil, errors.New("Can't Find Account")
 	}
	if account.CurrentBalance < req.Money {
		return nil, errors.New("InvalId Amount")
	}
	log.Println("Account Id: ", req.Id, " Withdraw Money: ", req.Money)
	trans_res := DB.Transaction(func(tx *gorm.DB) error {
		new_transaction := Transaction{AccountId : req.Id, Type : "Withdraw", Amount : req.Money, Date : time.Now()}
		if err := tx.Create(&new_transaction).Error; err != nil {
   			return err
 		}
		account.CurrentBalance -= req.Money
		if err := tx.Save(&account).Error; err != nil {
   			return err
 		}
		return nil
	})
	if trans_res != nil {
		return nil, errors.New("Can't make transaction")
	}
	log.Println("Account Id: ", req.Id, " Withdraw Money: ", req.Money, " Money after withdraw: ", account.CurrentBalance)
	response := &pb.WithdrawMoneyResponse{
		Result: "Withdraw success",
	}
	return response, nil
}

func (bs *BankingServer) BankAccountReport(ctx context.Context, req *pb.BankAccountReportRequest) (*pb.BankAccountReportResponse, error){
	var transactions []Transaction
	if err := DB.Where(&Transaction{ AccountId : req.AccountId}).Find(&transactions).Error; err!=nil{
		return nil, errors.New("Can't get Account report")
	}
	var total_deposit int64 = 0
	var total_withdraw int64 = 0
	trans_info := []string{}
	for _,transaction := range transactions {
		if transaction.Type == "Deposit" {
			total_deposit += transaction.Amount
		} else {
			total_withdraw += transaction.Amount
		}
		trans_info = append(trans_info, transaction.Type + ": " + strconv.FormatInt(transaction.Amount,10) + " Date: " + transaction.Date.String())
	}
	response := &pb.BankAccountReportResponse{
		AccountId: req.AccountId,
		Transactions: trans_info,
		TotalDeposit: total_deposit,
		TotalWithdraw: total_withdraw,
	}
	return response, nil
}
func (bs *BankingServer) AllAccountReport(ctx context.Context, req *pb.EmptyRequest) (*pb.AllAccountReportResponse, error) {
	accounts := []BankAccount{}
	if err := DB.Find(&accounts).Error; err != nil{
		return nil, errors.New("Can't get report")
	}
	reports := []*pb.BankAccountReportResponse{}
	
	for _,acc := range accounts {
		report, err := bs.BankAccountReport(ctx, &pb.BankAccountReportRequest{AccountId: acc.Id}); 
		if err != nil {
			return nil, errors.New("Can't get report")
		}
		reports = append(reports, report)
	}
	response := &pb.AllAccountReportResponse{
		BankAccountReport: reports,
	}
	return response, nil
}

func (bs *BankingServer) GetUserAllAccount(ctx context.Context, req *pb.GetUserAllAccountRequest) (*pb.GetUserAllAccountResponse, error) {
	accounts := []BankAccount{}
	if err := DB.Where(&BankAccount{ UserId : req.Id}).Find(&accounts).Error; err!=nil{
		return nil, errors.New("Can't delete account")
	}
	ids := []int32{}
	for _,acc := range accounts {
		ids = append(ids, acc.Id)
	}
	response := &pb.GetUserAllAccountResponse{
		BankAccountIds: ids,
	}
	return response, nil
}

func (bs *BankingServer) DeleteBankAccount(ctx context.Context, req *pb.DeleteBankAccountRequest) (*pb.DeleteBankAccountResponse, error) {

	trans_res := DB.Transaction(func(tx *gorm.DB) error {
		transactions := Transaction{}
		if err := tx.Where(&Transaction{ AccountId : req.Id}).Delete(&transactions).Error; err != nil {
   			return err
 		}
		account := BankAccount{}
		if err := tx.Where(&BankAccount{ Id : req.Id}).Delete(&account).Error; err!=nil{
			return err
		}
		return nil
	})
	if trans_res != nil {
		return nil, errors.New("Can't delete")
	}


	response := &pb.DeleteBankAccountResponse{
		Result : "Account deleted",
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
	// pb.RegisterUserManagementServer(s, &UserManagement{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
