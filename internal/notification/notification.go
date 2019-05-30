// Package notification handles sending notifications. At present it just wraps
// the 8o8/email package so that calles are agnostic as to how the
// notificfations are implemented.
package notification

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/8o8/email"
)

// Attachment is a copy of email.Attachment
type Attachment struct {
	MIMEType      string
	FileName      string
	Base64Content string
}

// Email is a copy of email.Email
type Email struct {
	FromName     string
	FromEmail    string
	ToName       string
	ToEmail      string
	Subject      string
	PlainContent string
	HTMLContent  string
	Attachments  []Attachment
}

// Send sends an email using the default mx service specified in the .env
func (e Email) Send() error {

	// map local Email value to email.Email
	eml := email.Email{
		FromName:     e.FromName,
		FromEmail:    e.FromEmail,
		ToName:       e.ToName,
		ToEmail:      e.ToEmail,
		Subject:      e.Subject,
		PlainContent: e.PlainContent,
		HTMLContent:  e.HTMLContent,
	}
	for _, a := range e.Attachments {
		att := email.Attachment{
			MIMEType:      a.MIMEType,
			FileName:      a.FileName,
			Base64Content: a.Base64Content,
		}
		eml.Attachments = append(eml.Attachments, att)
	}

	// get the preferred mx service from the env
	mx := os.Getenv("MAPPCPD_MX_SERVICE")
	if mx == "" {
		return errors.New("notification.Send() could not get the preferred MX service from env var MAPPCPD_MX_SERVICE")
	}

	switch strings.ToLower(mx) {
	case "mailgun":
		log.Printf("Sending email to %s via Mailgun", eml.ToEmail)
		return sendMailgun(
			eml,
			os.Getenv("MAILGUN_API_KEY"),
			os.Getenv("MAILGUN_DOMAIN"),
		)
	case "sendgrid":
		log.Printf("Sending email to %s via Sendgrid", eml.ToEmail)
		return sendSendgrid(
			eml,
			os.Getenv("SENDGRID_API_KEY"),
		)
	case "ses":
		log.Printf("Sending email to %s via SES", eml.ToEmail)
		return sendSES(
			eml,
			os.Getenv("AWS_SES_REGION"),
			os.Getenv("AWS_SES_ACCESS_KEY_ID"),
			os.Getenv("AWS_SES_SECRET_ACCESS_KEY"),
		)
	}

	return fmt.Errorf("notification.Send() unknown value for MAPPCPD_MX_SERVICE %q", mx)
}

// sendSES sends the email with Amazon SES
func sendSES(eml email.Email, awsRegion, awsAccessKeyID, awsSecretAccessKey string) error {
	cfg := email.SESCfg{
		AWSRegion:          awsRegion,
		AWSAccessKeyID:     awsAccessKeyID,
		AWSSecretAccessKey: awsSecretAccessKey,
	}
	sndr, err := email.NewSES(cfg)
	if err != nil {
		return err
	}

	return sndr.Send(eml)
}

// sendMailgun sends the email using Mailgun
func sendMailgun(eml email.Email, apiKey, domain string) error {
	cfg := email.MailgunCfg{
		APIKey: apiKey,
		Domain: domain,
	}
	sndr := email.NewMailgun(cfg)

	return sndr.Send(eml)
}

// sendSendgrid sends the email using Sendgrid
func sendSendgrid(eml email.Email, apiKey string) error {
	cfg := email.SendgridCfg{
		APIKey: apiKey,
	}
	sndr := email.NewSendgrid(cfg)

	return sndr.Send(eml)
}
