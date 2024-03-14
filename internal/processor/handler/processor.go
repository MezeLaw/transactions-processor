package handler

import (
	"context"
	"errors"
	"log"
	"transactions-processor/internal/processor/models"
)

const (
	s3Error                        = "an error occurred trying to retrieve file from S3"
	processingError                = "an error occurred when processing transactions"
	saveResultErr                  = "an error occurred trying to save results"
	transactionProcessingCompleted = "the transaction processing was completed"
)

type ProcessorService interface {
	GetTransactionsFromFile() ([]*models.Transaction, error)
	Process(transactions []*models.Transaction) (*models.TransactionResult, error)
	SaveResult(result *models.TransactionResult) error
}

type Processor struct {
	service ProcessorService
}

func New(processorService ProcessorService) *Processor {
	return &Processor{service: processorService}
}

func (h *Processor) Handle(ctx context.Context) (string, error) {
	log.Println("Starting transaction processing flow")
	transactions, err := h.service.GetTransactionsFromFile()

	if err != nil {
		log.Println(s3Error)
		return "S3 Error", errors.New(s3Error)
	}

	result, err := h.service.Process(transactions)

	println("Balance> ", result.Balance)
	println("Average Debit", result.ExtraInformation.AverageDebit)
	println("Average Credit", result.ExtraInformation.AverageCredit)
	println("September Qty", result.ExtraInformation.MonthlyTransactions["September"])
	println("October Qty", result.ExtraInformation.MonthlyTransactions["October"])

	if err != nil {
		log.Println(processingError)
		return "Processing Error", errors.New(processingError)
	}

	if saveErr := h.service.SaveResult(result); saveErr != nil {
		log.Println(saveResultErr)
		return "Save Result", errors.New(saveResultErr)
	}

	log.Println("Transaction processing completed")

	return transactionProcessingCompleted, nil

}
