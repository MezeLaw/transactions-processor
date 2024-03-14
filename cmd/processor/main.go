package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"transactions-processor/internal/processor/handler"
	"transactions-processor/internal/processor/service"
)

func main() {
	s := service.New(nil)
	h := handler.New(s)
	lambda.Start(h.Handle)

}
