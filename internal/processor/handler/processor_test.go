package handler

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"transactions-processor/internal/processor/models"
)

func TestProcessor_Handle(t *testing.T) {

	s3Transactions := []*models.Transaction{{
		ID:     "acc-01",
		Date:   "9-12",
		Amount: 77,
		Type:   "CREDIT",
	}}

	transactionResult := &models.TransactionResult{
		AccountID: "acc-01",
		Balance:   "77",
		ExtraInformation: models.ExtraInformation{
			MonthlyTransactions: map[string]int{
				"September": 1,
			},
			AverageDebit:  "77",
			AverageCredit: "0",
		},
	}

	tests := []struct {
		name           string
		result         string
		error          error
		mockedBehavior func(mock *mock.Mock)
		asserts        func(t *testing.T, expectedResult, result string, expectedError, error error)
	}{
		{
			name:   "execution should process successfully the transactions",
			result: "the transaction processing was completed",
			error:  nil,
			mockedBehavior: func(mock *mock.Mock) {
				mock.On("GetTransactionsFromFile").Return(s3Transactions, nil)
				mock.On("Process", s3Transactions).Return(transactionResult, nil)
				mock.On("SaveResult", transactionResult).Return(nil)
			},
			asserts: func(t *testing.T, expectedResult, result string, expectedError, error error) {
				assert.Equal(t, expectedResult, result)
				assert.Equal(t, expectedError, error)
			},
		},
		{
			name:   "execution should return error when saving results on dynamodb",
			result: "Save Result",
			error:  errors.New("an error occurred trying to save results"),
			mockedBehavior: func(mock *mock.Mock) {
				mock.On("GetTransactionsFromFile").Return(s3Transactions, nil)
				mock.On("Process", s3Transactions).Return(transactionResult, nil)
				mock.On("SaveResult", transactionResult).Return(errors.New("mocked err"))
			},
			asserts: func(t *testing.T, expectedResult, result string, expectedError, error error) {
				assert.Equal(t, expectedResult, result)
				assert.Equal(t, expectedError, error)
			},
		},
		{
			name:   "execution should return error when trying to execute process",
			result: "Processing Error",
			error:  errors.New("an error occurred when processing transactions"),
			mockedBehavior: func(mock *mock.Mock) {
				mock.On("GetTransactionsFromFile").Return(s3Transactions, nil)
				mock.On("Process", s3Transactions).Return(nil, errors.New("error"))
			},
			asserts: func(t *testing.T, expectedResult, result string, expectedError, error error) {
				assert.Equal(t, expectedResult, result)
				assert.Equal(t, expectedError, error)
			},
		},
		{
			name:   "execution should return error when trying to retrieve transactions from file",
			result: "S3 Error",
			error:  errors.New("an error occurred trying to retrieve file from S3"),
			mockedBehavior: func(mock *mock.Mock) {
				mock.On("GetTransactionsFromFile").Return(nil, errors.New("error"))
			},
			asserts: func(t *testing.T, expectedResult, result string, expectedError, error error) {
				assert.Equal(t, expectedResult, result)
				assert.Equal(t, expectedError, error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedService := ProcessorServiceMock{}
			handler := New(&mockedService)

			tt.mockedBehavior(&mockedService.Mock)

			result, err := handler.Handle(context.TODO())

			tt.asserts(t, tt.result, result, tt.error, err)

		})
	}

}

type ProcessorServiceMock struct {
	mock.Mock
}

func (m *ProcessorServiceMock) GetTransactionsFromFile() ([]*models.Transaction, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Transaction), nil
}
func (m *ProcessorServiceMock) Process(transactions []*models.Transaction) (*models.TransactionResult, error) {
	args := m.Called(transactions)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.TransactionResult), nil
}
func (m *ProcessorServiceMock) SaveResult(result *models.TransactionResult) error {
	args := m.Called(result)
	return args.Error(0)

}
