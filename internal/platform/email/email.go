// Package email sends emails using Sendgrid API. At this stage can only do single emails with one attachemnt.
package email

import (
	"log"
	"os"

	"github.com/34South/envr"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Email represents an email
type Email struct {
	FromName     string
	FromEmail    string
	ToEmail      string
	ToName       string
	Subject      string
	PlainContent string
	HTMLContent  string
	Attachments  []Attachment
}

// Attachment to an email
type Attachment struct {
	MIMEType      string
	FileName      string
	Base64Content string
}

// New return a pointer to an Email.
func New() *Email {
	return &Email{}
}

func init() {
	envr.New("email", []string{"SENDGRID_API_KEY"}).Clean()
}

// Send sends an email
func (e Email) Send() error {

	message := prepare(e)

	for _, a := range e.Attachments {
		attach(a, message)
	}

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		log.Println(err)
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}

	return err
}

func (a Attachment) isOK() bool {
	return a.Base64Content != "" && a.FileName != "" && a.MIMEType != ""
}

func prepare(e Email) *mail.SGMailV3 {
	from := mail.NewEmail(e.FromName, e.FromEmail)
	subject := e.Subject
	to := mail.NewEmail(e.ToName, e.ToEmail)
	plainTextContent := e.PlainContent
	htmlContent := e.HTMLContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	return message
}

func attach(a Attachment, message *mail.SGMailV3) {
	if a.isOK() {
		message.AddAttachment(newAttachment(a))
	}
}

func newAttachment(a Attachment) *mail.Attachment {
	ma := mail.NewAttachment()
	ma.SetContent(a.Base64Content)
	ma.SetType(a.MIMEType)
	ma.SetFilename(a.FileName)
	ma.SetDisposition("attachment") // no "inline" for now
	// ma.SetContentID("Attachment...") // used for inline attachments
	return ma
}
