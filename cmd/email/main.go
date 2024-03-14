package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	emailHandler "transactions-processor/internal/email/handler"
)

func main() {

	h := emailHandler.New()
	lambda.Start(h.SendEmail)
}
