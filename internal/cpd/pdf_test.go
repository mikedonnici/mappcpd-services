package cpd_test

import (
	"os"
	"testing"

	"github.com/cardiacsociety/web-services/internal/cpd"
	"github.com/matryer/is"
)

func TestCreatePDF(t *testing.T) {
	is := is.New(t)
	f, err := os.Create(os.TempDir() + "/test.pdf")
	defer f.Close()
	is.NoErr(err) // Error creating pdf file

	m := cpd.MemberActivityReport{
		ID: 1,
	}

	err = cpd.PDFReport(m, f)
	is.NoErr(err) // Could not create PDF
}
