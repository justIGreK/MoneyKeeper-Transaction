package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/justIGreK/MoneyKeeper-Transaction/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepo struct {
	collection *mongo.Collection
}

func NewTransactionRepository(db *mongo.Client) *TransactionRepo {
	return &TransactionRepo{
		collection: db.Database(dbname).Collection(transactionCollection),
	}
}

func (r *TransactionRepo) AddTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	result, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *TransactionRepo) GetTransaction(ctx context.Context, transactionID, userID string) (*models.Transaction, error) {
	oid, err := convertToObjectIDs(transactionID)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	log.Println(oid, userID)
	var transaction models.Transaction
	err = r.collection.FindOne(ctx, bson.M{"_id": oid[0], "user_id": userID}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, err
}

func (r *TransactionRepo) GetTXByTimeFrame(ctx context.Context, userID string, dateFrame models.TimeFrame) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gt": dateFrame.StartDate,
			"$lt": dateFrame.EndDate,
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, err
}

func (r *TransactionRepo) GetAllTransactions(ctx context.Context, userID string) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, err
}
