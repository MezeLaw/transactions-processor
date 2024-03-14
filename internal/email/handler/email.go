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

			if err != nil {
				log.Println("an error occurred trying to marshall received item info")
				return err
			}
			balance := accountInfo["balance"].String()
			averageDebit := accountInfo["average_debit"].String()
			averageCredit := accountInfo["average_credit"].String()
			jan := accountInfo["monthly_transactions"].Map()["January"].String()
			feb := accountInfo["monthly_transactions"].Map()["February"].String()
			mar := accountInfo["monthly_transactions"].Map()["March"].String()
			apr := accountInfo["monthly_transactions"].Map()["April"].String()
			may := accountInfo["monthly_transactions"].Map()["May"].String()
			jun := accountInfo["monthly_transactions"].Map()["June"].String()
			jul := accountInfo["monthly_transactions"].Map()["July"].String()
			aug := accountInfo["monthly_transactions"].Map()["August"].String()
			sept := accountInfo["monthly_transactions"].Map()["September"].String()
			oct := accountInfo["monthly_transactions"].Map()["October"].String()
			nov := accountInfo["monthly_transactions"].Map()["November"].String()
			dec := accountInfo["monthly_transactions"].Map()["December"].String()

			messageInfo := `
							Total balance is: ` + balance + `
							Average debit amount: ` + averageDebit + `
							Average credit amount: ` + averageCredit + `

							Number of transactions in January: ` + jan + `
							Number of transactions in February: ` + feb + `
							Number of transactions in March: ` + mar + `	
							Number of transactions in April: ` + apr + `
							Number of transactions in May: ` + may + `
							Number of transactions in June: ` + jun + `
							Number of transactions in July: ` + jul + `
							Number of transactions in August:` + aug + `
							Number of transactions in September: ` + sept + `
							Number of transactions in October: ` + oct + `
							Number of transactions in November: ` + nov + `
							Number of transactions in December: ` + dec + `
							`

			input := &ses.SendEmailInput{
				Destination: &ses.Destination{
					ToAddresses: []*string{
						aws.String(defaultEmailAddress),
					},
				},
				Message: &ses.Message{
					Body: &ses.Body{
						Text: &ses.Content{
							Data: aws.String(messageInfo),
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

			fmt.Println("Email was sent successfully")
		}
	}
	return nil
}
