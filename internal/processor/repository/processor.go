package repository

import "transactions-processor/internal/processor/models"

type Processor struct {
}

func New() Processor {
	return Processor{}
}

func (r *Processor) Save(*models.TransactionResult) error {
	return nil
}
