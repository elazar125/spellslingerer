package mailxtns

import (
	"fmt"
	"log"
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func BindMailEvents(app core.App) {
	app.OnMailerBeforeAdminResetPasswordSend().Add(func(e *core.MailerAdminEvent) error {
		if err := sendMail(*e.Message); err != nil {
			return err
		}
		return hook.StopPropagation
	})
	app.OnMailerBeforeRecordResetPasswordSend().Add(func(e *core.MailerRecordEvent) error {
		if err := sendMail(*e.Message); err != nil {
			return err
		}
		return hook.StopPropagation
	})
	app.OnMailerBeforeRecordChangeEmailSend().Add(func(e *core.MailerRecordEvent) error {
		if err := sendMail(*e.Message); err != nil {
			return err
		}
		return hook.StopPropagation
	})
	app.OnMailerBeforeRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		if err := sendMail(*e.Message); err != nil {
			return err
		}
		return hook.StopPropagation
	})
}

func sendMail(pbMessage mailer.Message) error {
	to := mail.NewEmail(pbMessage.To[0].Name, pbMessage.To[0].Address)
	from := mail.NewEmail(pbMessage.From.Name, pbMessage.From.Address)
	subject := pbMessage.Subject
	htmlContent := pbMessage.HTML
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SPELLSLINGERER_SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return err
}
