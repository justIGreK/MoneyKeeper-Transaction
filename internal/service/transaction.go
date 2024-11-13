package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/justIGreK/MoneyKeeper-Transaction/internal/models"
)

type TransactionRepository interface {
	AddTransaction(ctx context.Context, transaction models.Transaction) (string, error)
	GetTransaction(ctx context.Context, transactionID, userID string) (*models.Transaction, error)
	GetAllTransactions(ctx context.Context, userID string) ([]models.Transaction, error)
	GetTXByTimeFrame(ctx context.Context, userID string, dateFrame models.TimeFrame) ([]models.Transaction, error)
	UpdateTx(ctx context.Context, updates models.Transaction) error
	DeleteTx(ctx context.Context, userID, txID string) error
}

type TransactionService struct {
	TransactionRepo TransactionRepository
	User            UserService
}

type UserService interface {
	GetUser(ctx context.Context, id string) (string, string, error)
}

func NewTransactionService(transRepo TransactionRepository, user UserService) *TransactionService {
	return &TransactionService{TransactionRepo: transRepo,
		User: user}
}

const (
	NoCategory            = "other"
	Dateformat     string = "2006-01-02"
	DateTimeformat string = "2006-01-02T15:04:05"
	TimeFormat     string = "15:04"
)

func (s *TransactionService) AddTransaction(ctx context.Context, transaction models.CreateTransaction) (string, error) {
	if transaction.Cost < 0 {
		return "", errors.New("cost cant be below 0")
	}
	id, _, err := s.User.GetUser(ctx, transaction.UserID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if id == "" {
		return "", errors.New("user not found")
	}
	if transaction.Category == "" {
		transaction.Category = NoCategory
	}
	now := time.Now().UTC()
	date := now
	if transaction.Date != nil {
		date, err = time.Parse(DateTimeformat, *transaction.Date)
		if err != nil {
			return "", err
		}
	}
	createTransaction := models.Transaction{
		UserID:   transaction.UserID,
		Category: transaction.Category,
		Name:     transaction.Name,
		Cost:     transaction.Cost,
		Date:     date,
	}
	id, err = s.TransactionRepo.AddTransaction(ctx, createTransaction)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, transactionID, userID string) (*models.Transaction, error) {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	trans, err := s.TransactionRepo.GetTransaction(ctx, transactionID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return trans, nil
}

func (s *TransactionService) GetAllTransactions(ctx context.Context, userID string) ([]models.Transaction, error) {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}

	txs, err := s.TransactionRepo.GetAllTransactions(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return txs, nil
}
func (s *TransactionService) GetTXByTimeFrame(ctx context.Context, userID string, timeframe models.CreateTimeFrame) ([]models.Transaction, error) {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	var tf models.TimeFrame
	if timeframe.StartDate == "" {
		tf.StartDate = time.Unix(0, 0)
	} else {
		date, err := time.Parse(Dateformat, timeframe.StartDate)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		tf.StartDate = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 1, date.Location())
	}
	if timeframe.EndDate == "" {
		tf.EndDate = time.Now().AddDate(10000, 0, 0)
	} else {
		date, err := time.Parse(Dateformat, timeframe.EndDate)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		tf.EndDate = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 9999999, date.Location())
	}

	txs, err := s.TransactionRepo.GetTXByTimeFrame(ctx, userID, tf)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return txs, nil
}
func (s *TransactionService) DeleteTx(ctx context.Context, userID, txID string) error {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	if user == "" {
		return errors.New("user not found")
	}

	tx, err := s.TransactionRepo.GetTransaction(ctx, txID, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	if tx == nil{
		return errors.New("transaction is not found")
	}
	err = s.TransactionRepo.DeleteTx(ctx, txID, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (s *TransactionService) UpdateTx(ctx context.Context, updates models.UpdateTransaction) (*models.Transaction, error) {
	user, _, err := s.User.GetUser(ctx, updates.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	tx, err := s.TransactionRepo.GetTransaction(ctx, updates.ID, updates.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("transaction is not found")
	}
	updatedTx := models.Transaction{
		ID:     updates.ID,
		UserID: updates.UserID,
	}
	if updates.Name != nil {
		updatedTx.Name = *updates.Name
	} else {
		updatedTx.Name = tx.Name
	}
	if updates.Cost != nil {
		updatedTx.Cost = *updates.Cost
	} else {
		updatedTx.Cost = tx.Cost
	}
	if updates.Category != nil {
		updatedTx.Category = *updates.Category
	} else {
		updatedTx.Category = tx.Category
	}
	updatedTx.Date, err = s.parseDateTime(tx, updates)

	err = s.TransactionRepo.UpdateTx(ctx, updatedTx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	newTx, err := s.TransactionRepo.GetTransaction(ctx, updates.ID, updates.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return newTx, nil
}

func (s *TransactionService) parseDateTime(tx *models.Transaction, updates models.UpdateTransaction) (time.Time, error) {
	var date, times time.Time
	var err error
	if updates.Date == nil && updates.Time == nil {
		return tx.Date, nil
	}
	if updates.Date != nil {
		date, err = time.Parse(Dateformat, *updates.Date)
		if err != nil{
			return time.Time{}, err
		}
	}else {
		date = tx.Date
	}
	if updates.Time != nil {
		times, err = time.Parse(TimeFormat, *updates.Time)
		if err != nil{
			return time.Time{}, err
		}
	}else {
		times = tx.Date
	} 
	dateResp := time.Date(date.Year(), date.Month(), date.Day(), times.Hour(), times.Minute(), times.Second(), 0, time.Now().UTC().Location())
	return dateResp, nil
}
