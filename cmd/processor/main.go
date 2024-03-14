package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"transactions-processor/internal/processor/handler"
	"transactions-processor/internal/processor/repository"
	"transactions-processor/internal/processor/service"
)

func main() {
	r := repository.New()
	s := service.New(&r)
	h := handler.New(s)
	lambda.Start(h.Handle)

}
