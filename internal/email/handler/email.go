package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

const (
	newAccount          = "INSERT"
	updatedAccount      = "MODIFY"
	defaultEmailAddress = "mezequielabogado@gmail.com" // TODO para mejorarlo puede agregarse en la lambda inicial en el request body el email y al insertarlo levantarlo y utilizarlo para el envio de email
	emailSubject        = "Account Balance - Info "
	senderEmailAddress  = "mezequielabogado@gmail.com"
)

type Email struct{}

func New() Email {
	return Email{}
}

func (e *Email) SendEmail(ctx context.Context, event events.DynamoDBEvent) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svcSES := ses.New(sess)

	for _, record := range event.Records {
		if record.EventName == newAccount || record.EventName == updatedAccount {

			accountInfo := record.Change.NewImage

			//accountInfoJson, err := json.Marshal(accountInfo)

			if err != nil {
				log.Println("an error occurred trying to marshall received item info")
				return err
			}
			balance := accountInfo["balance"].String()
			averageDebit := accountInfo["average_debit"].String()
			averageCredit := accountInfo["average_credit"].String()

			messageInfo := `	Transactions info 
								Total balance is: ` + balance + `
								Average debit amount: ` + averageDebit + `
								Average credit amount: ` + averageCredit + `
								Number of transactions in January: ` + "jan" + `
								Number of transactions in February: ` + "feb" + `
								Number of transactions in March: ` + "mar" + `
								Number of transactions in April: ` + "apr" + `
								Number of transactions in May: ` + "may" + `
								Number of transactions in June: ` + "jun" + `
								Number of transactions in July: ` + "jul" + `
								Number of transactions in August:` + "aug" + `
								Number of transactions in September: ` + "sept" + `
								Number of transactions in October: ` + "oct" + `
								Number of transactions in November: ` + "nov" + `
								Number of transactions in December: ` + "dec" + `
								<img src="https://trx-public-bucket.s3.amazonaws.com/download.png" alt="Stori Logo">
							`
			// Crear cuerpo del correo electrónico
			bodyText := fmt.Sprintf("%s", messageInfo)

			input := &ses.SendEmailInput{
				Destination: &ses.Destination{
					ToAddresses: []*string{
						aws.String(defaultEmailAddress),
					},
				},
				Message: &ses.Message{
					Body: &ses.Body{
						Text: &ses.Content{
							Data: aws.String(bodyText),
						},
					},
					Subject: &ses.Content{
						Data: aws.String(emailSubject),
					},
				},
				Source: aws.String(senderEmailAddress),
			}

			_, err = svcSES.SendEmail(input)
			if err != nil {
				log.Println("An error occurred trying to send email. Err: ", err.Error())
				return err
			}

			fmt.Println("Correo electrónico enviado correctamente")
		}
	}
	return nil
}
