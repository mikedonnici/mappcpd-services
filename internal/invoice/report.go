package invoice

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/platform/excel"
)

// ExcelReport returns an excel invoice report File
func ExcelReport(ds datastore.Datastore, invoices []Invoice) (*excelize.File, error) {

	f := excel.New([]string{
		"Invoice ID",
		"Invoice date",
		"Due date",
		"Subscription",
		"Amount",
		"Paid",
		"Comment",
		"Member ID",
		"Name",
		"Email",
		"Mobile",
		"Entry date",
		"Membership",
		"Status",
		"Country",
		"Tags",
		"Journal num.",
		"BPAY num.",
		"Address",
		"Locality",
		"State",
		"Postcode",
		"Country",
	})

	// data rows
	var total float64
	for _, i := range invoices {

		paid := "no"
		if i.Paid == true {
			paid = "yes"
		}

		data := []interface{}{
			i.ID,
			i.IssueDate,
			i.DueDate,
			i.Subscription,
			i.Amount,
			paid,
			i.Comment,
			i.MemberID,
			i.Member.Title + " " + i.Member.FirstName + " " + i.Member.LastName,
			i.Member.Contact.EmailPrimary,
			i.Member.Contact.Mobile,
			i.Member.DateOfEntry,
			i.Member.Memberships[0].Title,
			i.Member.Memberships[0].Status,
			i.Member.Country,
			strings.Join(i.Member.Tags, ", "),
			i.Member.JournalNumber,
			i.Member.BpayNumber,
			strings.Join(i.Member.Contact.Locations[0].Address, " "),
			i.Member.Contact.Locations[0].City,
			i.Member.Contact.Locations[0].State,
			i.Member.Contact.Locations[0].Postcode,
			i.Member.Contact.Locations[0].Country,
		}
		err := f.AddRow(data)
		if err != nil {
			msg := fmt.Sprintf("AddRow() err = %s", err)
			log.Printf(msg)
			f.AddError(i.ID, msg)
		}

		total += i.Amount
	}

	// total row
	r := []interface{}{
		"", "", "", "Total", total,
		"", "", "", "", "", "", "",
		"", "", "", "", "", "", "",
		"", "", "", "",
	}
	err := f.AddRow(r)
	if err != nil {
		msg := fmt.Sprintf("AddRow() err = %s\n", err)
		log.Printf(msg)
		f.AddError(0, msg)
	}

	// style
	f.SetColStyleByHeading("Invoice date", excel.DateStyle)
	f.SetColWidthByHeading("Invoice date", 18)
	f.SetColStyleByHeading("Due date", excel.DateStyle)
	f.SetColWidthByHeading("Due date", 18)
	f.SetColWidthByHeading("Name", 18)
	f.SetColStyleByHeading("Amount", excel.CurrencyStyle)
	f.SetColWidthByHeading("Amount", 18)
	cell := "E" + strconv.Itoa(f.NextRow)
	f.SetCellStyle(cell, cell, excel.BoldStyle)
	cell = "F" + strconv.Itoa(f.NextRow)
	f.SetCellStyle(cell, cell, excel.BoldCurrencyStyle)

	return f.XLSX, nil
}
