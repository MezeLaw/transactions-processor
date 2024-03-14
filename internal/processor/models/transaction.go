package models

type Transaction struct {
	ID     string
	Date   string
	Amount float64
	Type   string
}

type TransactionResult struct {
	AccountID        string
	Balance          string
	ExtraInformation ExtraInformation
}

type ExtraInformation struct {
	MonthlyTransactions map[string]int
	AverageDebit        string
	AverageCredit       string
}
