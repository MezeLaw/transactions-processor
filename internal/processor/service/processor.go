package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"strconv"
	"strings"
	"transactions-processor/internal/processor/models"
)

const (
	defaultAccountID = "acc-01"
	debitType        = "DEBIT"
	creditType       = "CREDIT"
	january          = "January"
	february         = "February"
	march            = "March"
	april            = "April"
	may              = "May"
	june             = "June"
	july             = "July"
	august           = "August"
	september        = "September"
	october          = "October"
	november         = "November"
	december         = "December"
)

type ProcessorRepository interface {
	Save(*models.TransactionResult) error
}

type Processor struct {
	repository ProcessorRepository
}

func New(repository ProcessorRepository) *Processor {
	return &Processor{repository: repository}
}

func (s *Processor) GetTransactionsFromFile() ([]*models.Transaction, error) {

	file, err := s.getCSVFromS3()
	if err != nil {
		return nil, err
	}

	trx, err := s.getTransactionsFromFile(file)

	return trx, nil
}

func (s *Processor) Process(transactions []*models.Transaction) (*models.TransactionResult, error) {

	balance := s.getBalance(transactions)
	extraInfo := s.getExtraInfo(transactions)
	return &models.TransactionResult{
		AccountID:        defaultAccountID,
		Balance:          fmt.Sprintf("%.2f", balance),
		ExtraInformation: extraInfo,
	}, nil
}

func (s *Processor) SaveResult(result *models.TransactionResult) error {

	err := s.repository.Save(result)
	if err != nil {
		println("error trying to save account")
		return err
	}

	return nil
}

func (s *Processor) getCSVFromS3() (*s3.GetObjectOutput, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	bucket := "trx-public-bucket"
	key := "transactions.csv"

	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		fmt.Println("Error trying to retrieve object from S3:", err)
		return nil, errors.New("error trying to retrieve object from S3")
	}
	return resp, nil
}

func (s *Processor) getTransactionsFromFile(file *s3.GetObjectOutput) ([]*models.Transaction, error) {

	transactions := []*models.Transaction{}

	csvReader := csv.NewReader(file.Body)

	_, err := csvReader.Read()
	if err != nil {
		fmt.Println("Error trying to read first line of CSV files:", err)
		return nil, err
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error trying to read CSV line:", err)
			return nil, err
		}

		amount, err := s.parseAmount(record[2])
		if err != nil {
			return nil, errors.New("error parsing trx amount")
		}

		transaction := &models.Transaction{
			ID:     record[0],
			Date:   record[1],
			Amount: *amount,
			Type:   s.setTransactionType(record[2]),
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *Processor) setTransactionType(amount string) string {
	if strings.Contains(amount, "-") {
		return debitType
	}
	return creditType
}

func (s *Processor) parseAmount(amount string) (*float64, error) {

	floatValue, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println("Error trying to parse amount:", err)
		return nil, errors.New("error parsing amounts")
	}

	return &floatValue, nil
}

func (s *Processor) getBalance(transactions []*models.Transaction) float64 {

	balance := 0.00

	for _, trx := range transactions {
		balance += trx.Amount
	}

	return balance
}

func (s *Processor) getExtraInfo(transactions []*models.Transaction) models.ExtraInformation {

	trxByMonth := map[string]int{
		january:   0,
		february:  0,
		march:     0,
		april:     0,
		may:       0,
		june:      0,
		july:      0,
		august:    0,
		september: 0,
		october:   0,
		november:  0,
		december:  0,
	}

	totalDebit := 0.00
	totalCredit := 0.00

	debitQty := 0
	creditQty := 0

	for _, trx := range transactions {
		if trx.Amount < 0 {
			totalDebit += trx.Amount
			debitQty += 1
			trxByMonth[s.getTransactionMonth(trx.Date)] = trxByMonth[s.getTransactionMonth(trx.Date)] + 1
		} else {
			totalCredit += trx.Amount
			creditQty += 1
			trxByMonth[s.getTransactionMonth(trx.Date)] = trxByMonth[s.getTransactionMonth(trx.Date)] + 1
		}
	}

	extraInfo := models.ExtraInformation{
		MonthlyTransactions: trxByMonth,
		AverageDebit:        fmt.Sprintf("%.2f", totalDebit/float64(debitQty)),
		AverageCredit:       fmt.Sprintf("%.2f", totalCredit/float64(creditQty)),
	}

	return extraInfo
}

func (s *Processor) getTransactionMonth(date string) string {
	dateSplit := strings.Split(date, "-")
	month := dateSplit[0]

	switch month {
	case "1":
		return january
	case "2":
		return february
	case "3":
		return march
	case "4":
		return april
	case "5":
		return may
	case "6":
		return june
	case "7":
		return july
	case "8":
		return august
	case "9":
		return september
	case "10":
		return october
	case "11":
		return november
	case "12":
		return december
	default:
		return september // lo elegi arbitrariamente yo
	}

}
