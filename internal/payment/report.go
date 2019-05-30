package payment

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/platform/excel"
)

// ExcelReport returns an excel payment report File
func ExcelReport(ds datastore.Datastore, payments []Payment) (*excelize.File, error) {

	f := excel.New([]string{
		"Payment ID",
		"Payment date",
		"Member",
		"Payment type",
		"Amount",
		"Invoice",
		"Comment",
	})

	// data rows
	var total float64
	for _, p := range payments {

		var ia []string
		for _, i := range p.Allocations {
			ia = append(ia, strconv.Itoa(i.InvoiceID))
		}
		invoiceAllocations := strings.Join(ia, ", ")

		data := []interface{}{
			p.ID,
			p.Date,
			p.Member + " [" + strconv.Itoa(p.MemberID) + "]",
			p.Type,
			p.Amount,
			invoiceAllocations,
			p.Comment,
		}
		err := f.AddRow(data)
		if err != nil {
			msg := fmt.Sprintf("AddRow() err = %s", err)
			log.Printf(msg)
			f.AddError(p.ID, msg)
		}

		total += p.Amount
	}

	// total row
	r := []interface{}{"", "", "", "Total", total, "", ""}
	err := f.AddRow(r)
	if err != nil {
		msg := fmt.Sprintf("AddRow() err = %s", err)
		log.Printf(msg)
		f.AddError(0, msg)
	}

	// style
	f.SetColStyleByHeading("Payment date", excel.DateStyle)
	f.SetColWidthByHeading("Payment date", 18)
	f.SetColWidthByHeading("Member", 18)
	f.SetColStyleByHeading("Amount", excel.CurrencyStyle)
	f.SetColWidthByHeading("Amount", 18)
	cell := "D" + strconv.Itoa(f.NextRow)
	f.SetCellStyle(cell, cell, excel.BoldStyle)
	cell = "E" + strconv.Itoa(f.NextRow)
	f.SetCellStyle(cell, cell, excel.BoldCurrencyStyle)

	return f.XLSX, nil
}
