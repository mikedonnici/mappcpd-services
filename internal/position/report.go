package position

import (
	"fmt"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/platform/excel"
)

// ExcelReport returns an excel position report File
func ExcelReport(ds datastore.Datastore, positions []Position) (*excelize.File, error) {

	f := excel.New([]string{
		"ID",
		"Member",
		"Email",
		"Position",
		"Organisation",
		"Start",
		"End",
		"Comment",
	})

	// data rows
	for _, p := range positions {

		// If dates are bung set to an empty string
		var startDate, endDate interface{}
		if p.StartDate.Year() > 1971 { // epoch + 1
			startDate = p.StartDate
		} else {
			startDate = ""
		}
		if p.EndDate.Year() > 1971 { // epoch + 1
			endDate = p.EndDate
		} else {
			endDate = ""
		}

		data := []interface{}{
			p.MemberPositionID,
			p.Member + " [" + strconv.Itoa(p.MemberID) + "]",
			p.Email,
			p.Name + " [" + strconv.Itoa(p.ID) + "]",
			p.OrganisationName + " [" + strconv.Itoa(p.OrganisationID) + "]",
			startDate,
			endDate,
			p.Comment,
		}
		err := f.AddRow(data)
		if err != nil {
			msg := fmt.Sprintf("AddRow() err = %s", err)
			log.Printf(msg)
			f.AddError(p.ID, msg)
		}
	}

	// style
	f.SetColWidthByHeading("Member", 18)
	f.SetColStyleByHeading("Start", excel.DateStyle)
	f.SetColWidthByHeading("Start", 18)
	f.SetColStyleByHeading("End", excel.DateStyle)
	f.SetColWidthByHeading("End", 18)

	return f.XLSX, nil
}
