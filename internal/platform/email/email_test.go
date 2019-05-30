package email_test

import (
	"os"
	"testing"

	"github.com/34South/envr"
	"github.com/cardiacsociety/web-services/internal/platform/email"
)

func init() {
	envr.New("email", []string{"SENDGRID_API_KEY"}).Clean()
}

func TestSend(t *testing.T) {

	t.Log(os.Getenv("SENDGRID_API_KEY"))

	e := email.New()
	e.FromName = "Test Send"
	e.FromEmail = "test@test.com"
	e.ToName = "Mike Donnici"
	e.ToEmail = "michael@mesa.net.au"
	e.Subject = "Test email"
	e.HTMLContent = "<h1>Test</h1>"
	e.PlainContent = "test"
	if err := e.Send(); err != nil {
		t.Errorf("email.Send() err = %s", err)
	}
}
