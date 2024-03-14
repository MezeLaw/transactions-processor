package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strconv"
	"transactions-processor/internal/processor/models"
)

type Account struct {
	ID                  string            `json:"id"`
	Balance             string            `json:"balance"`
	MonthlyTransactions map[string]string `json:"monthly_transactions"`
	AverageDebit        string            `json:"average_debit"`
	AverageCredit       string            `json:"average_credit"`
}

type Processor struct {
}

func New() Processor {
	return Processor{}
}

func (r *Processor) Save(transactionResult *models.TransactionResult) error {
	// TODO mover la creacion del cliente de dynamo como un atributo del struct e inyectarlo en el main file
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	item := Account{
		ID:                  transactionResult.AccountID,
		Balance:             transactionResult.Balance,
		MonthlyTransactions: r.mapResultMonthlyTransactionsToAccountMonthlyTransactions(transactionResult.ExtraInformation.MonthlyTransactions),
		AverageDebit:        transactionResult.ExtraInformation.AverageDebit,
		AverageCredit:       transactionResult.ExtraInformation.AverageCredit,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println("error trying to marshall account item: ", err.Error())
	}

	tableName := "accounts"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("error to insert account item on db: %s", err)
	}

	return nil
}

func (r *Processor) mapResultMonthlyTransactionsToAccountMonthlyTransactions(trxs map[string]int) map[string]string {

	transactionsByMonth := map[string]string{}

	for month, trxQty := range trxs {

		transactionsByMonth[month] = strconv.Itoa(trxQty)

	}

	return transactionsByMonth
}
