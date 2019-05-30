package member

import (
	"log"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cardiacsociety/web-services/internal/platform/excel"
)

// ExcelReport returns an excel member report File
func ExcelReport(members []Member) (*excelize.File, error) {

	f := excel.New([]string{
		"Member ID",
		"Prefix",
		"First Name",
		"Middle Name(s)",
		"Last Name",
		"Suffix",
		"Gender",
		"Date of birth",
		"Email (primary)",
		"Email (secondary)",
		"Mobile",
		"Date of entry",
		"Membership Title",
		"Membership Status",
		"Membership Country",
		"Tags",
		"Journal No.",
		"BPAY No.",
		"Mail Address",
		"Mail Locality",
		"Mail State",
		"Mail Postcode",
		"Mail Country",
		"Directory Address",
		"Directory Locality",
		"Directory State",
		"Directory Postcode",
		"Directory Country",
		"Directory Phone",
		"Directory Fax",
		"Directory Email",
		"Directory Web",
		"First Council",
		"Second Council",
		"Third Council",
		"First Speciality",
		"Second Speciality",
		"Third Speciality",
	})

	// data rows
	for _, m := range members {

		var dob interface{}
		d, err := time.Parse("2006-01-02", m.DateOfBirth)
		if err == nil {
			dob = d // time.Time will accept the dateStyle formatting
		} else {
			f.AddError(m.ID, "Error parsing date of birth: "+err.Error())
		}

		var doe interface{}
		de, err := time.Parse("2006-01-02", m.DateOfEntry)
		if err == nil {
			doe = de // time.Time will accept the dateStyle formatting
		} else {
			f.AddError(m.ID, "Error parsing date of entry: "+err.Error())
		}

		var title string
		var status string
		if len(m.Memberships) > 0 {
			title = m.Memberships[0].Title
			status = m.Memberships[0].Status
		} else {
			f.AddError(m.ID, "Could not determine memership title / status")
		}

		var tags string
		if len(m.Tags) > 0 {
			tags = strings.Join(m.Tags, ", ")
		}

		// ContactLocationByType returns an empty struct and an error if not found
		// so can ignore error and write an empty cell
		mail, _ := m.ContactLocationByDesc("mail")
		directory, _ := m.ContactLocationByDesc("directory")

		p1, _ := m.PositionByName("First Council Affiliation")
		p2, _ := m.PositionByName("Second Council Affiliation")
		p3, _ := m.PositionByName("Third Council Affiliation")

		// There can be many specialities, but generally up to 3 for the report
		// they *should* be returned in order of preference
		var s1, s2, s3 string
		if len(m.Specialities) > 0 {
			s1 = m.Specialities[0].Name
		}
		if len(m.Specialities) > 1 {
			s2 = m.Specialities[1].Name
		}
		if len(m.Specialities) > 2 {
			s3 = m.Specialities[2].Name
		}

		data := []interface{}{
			m.ID,
			m.Title,
			m.FirstName,
			strings.Join(m.MiddleNames, " "),
			m.LastName,
			m.PostNominal,
			m.Gender,
			dob,
			m.Contact.EmailPrimary,
			m.Contact.EmailSecondary,
			m.Contact.Mobile,
			doe,
			title,
			status,
			m.Country,
			tags,
			m.JournalNumber,
			m.BpayNumber,
			strings.Join(mail.Address, " "),
			mail.City,
			mail.State,
			mail.Postcode,
			mail.Country,
			strings.Join(directory.Address, " "),
			directory.City,
			directory.State,
			directory.Postcode,
			directory.Country,
			directory.Phone,
			directory.Fax,
			directory.Email,
			directory.URL,
			p1.OrgName,
			p2.OrgName,
			p3.OrgName,
			s1,
			s2,
			s3,
		}
		err = f.AddRow(data)
		if err != nil {
			log.Printf("AddRow() err = %s\n", err)
			f.AddError(m.ID, err.Error())
		}
	}

	f.SetColStyleByHeading("Date of birth", excel.DateStyle)
	f.SetColWidthByHeading("Date of birth", 18)
	f.SetColStyleByHeading("Date of entry", excel.DateStyle)
	f.SetColWidthByHeading("Date of entry", 18)

	return f.XLSX, nil
}

// ExcelReportJournal is a cut down member report
func ExcelReportJournal(members []Member) (*excelize.File, error) {

	f := excel.New([]string{
		"Member ID",
		"Member",
		"Membership",
		"Journal no.",
		"Address 1",
		"Address 2",
		"Address 3",
		"Locality",
		"State",
		"Postcode",
		"Country",
		"Email",
	})

	// data rows
	for _, m := range members {

		var title string
		if len(m.Memberships) > 0 {
			title = m.Memberships[0].Title
		}

		// ContactLocationByType returns an empty struct and an error if not found
		// so can ignore error and write an empty cell
		var address = []string{"", "", ""}
		mail, err := m.ContactLocationByDesc("mail")
		if err != nil {
			f.AddError(m.ID, "Error fetching mail address: "+err.Error())
		}
		if len(mail.Address) > 0 {
			address[0] = mail.Address[0]
		}
		if len(mail.Address) > 1 {
			address[1] = mail.Address[1]
		}
		if len(mail.Address) > 2 {
			address[2] = mail.Address[2]
		}

		data := []interface{}{
			m.ID,
			m.Title + " " + m.FirstName + " " + m.LastName,
			title,
			m.JournalNumber,
			address[0],
			address[1],
			address[2],
			mail.City,
			mail.State,
			mail.Postcode,
			mail.Country,
			m.Contact.EmailPrimary,
		}
		err = f.AddRow(data)
		if err != nil {
			log.Printf("AddRow() err = %s\n", err)
			f.AddError(m.ID, err.Error())
		}
	}

	f.SetColWidthByHeading("Member", 18)
	f.SetColWidthByHeading("Address 1", 18)
	f.SetColWidthByHeading("Address 2", 18)
	f.SetColWidthByHeading("Address 3", 18)
	f.SetColWidthByHeading("Locality", 18)
	f.SetColWidthByHeading("Email", 18)

	return f.XLSX, nil
}
