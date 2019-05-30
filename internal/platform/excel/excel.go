// Package excel provides a simple abstraction for the excelize package
package excel

import (
	"errors"
	"log"
	"math"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const defaultSheetName = "Sheet1"
const errorSheetName = "Errors"
const defaultHeadingRow = 1
const defaultHeadingStyle = `{"font": {"bold": true}}`

const BoldStyle = `{"font": {"bold": true}}`
const CurrencyStyle = `{"number_format": 169}`
const BoldCurrencyStyle = `{"number_format": 169, "font":{"bold":true}}`
const DateStyle = `{"custom_number_format": "dd mmm yyyy"}`

// File represents a single-sheet xlsx file
type File struct {
	SheetName  string
	Columns    []column
	NextRow    int
	NextErrRow int
	XLSX       *excelize.File
}

type column struct {
	Ref         string
	Style       string
	HeadingCell string
	Heading     string
}

// New returns a pointer to an excel.File with all columns initialised to defaults
func New(colNames []string) *File {
	f := File{}
	f.SheetName = defaultSheetName
	f.NextRow = defaultHeadingRow
	f.NextErrRow = 2 // row 1 is headings
	f.XLSX = excelize.NewFile()
	xc := columnRefs(len(colNames))
	for i := range colNames {
		c := column{
			Ref:         xc[i],
			HeadingCell: xc[i] + strconv.Itoa(f.NextRow), // "A1", "A2" etc
			Heading:     colNames[i],
		}
		f.Columns = append(f.Columns, c)
		f.XLSX.SetCellValue(f.SheetName, c.HeadingCell, c.Heading)
	}

	f.SetHeadingStyle(defaultHeadingStyle)

	return &f
}

// SetHeadingStyle sets the column heading style
func (f *File) SetHeadingStyle(style string) {
	startCell := f.Columns[0].HeadingCell
	endCell := f.Columns[len(f.Columns)-1].HeadingCell
	f.SetCellStyle(startCell, endCell, style)
}

// SetAllColWidths sets all the column widths to the specified width
func (f *File) SetAllColWidths(width int) {
	startCell := f.Columns[0].Ref
	endCell := f.Columns[len(f.Columns)-1].Ref
	f.XLSX.SetColWidth(f.SheetName, startCell, endCell, float64(width))
}

// SetColWidthByHeading sets the width for a single column specified by the column heading
func (f *File) SetColWidthByHeading(heading string, width int) {
	c, err := f.colByHeading(heading)
	if err != nil {
		log.Printf("SetColWidthByHeading() err = %s\n", err)
		return
	}
	f.SetColWidth(c.Ref, width)
}

// SetColStyleByHeading sets the style for a single column specified by the column heading
func (f *File) SetColStyleByHeading(heading string, style string) {
	c, err := f.colByHeading(heading)
	if err != nil {
		log.Printf("SetColStyleByHeading() err = %s\n", err)
		return
	}
	startCell := c.Ref + strconv.Itoa(defaultHeadingRow+1)
	endCell := c.Ref + strconv.Itoa(f.NextRow)
	f.SetCellStyle(startCell, endCell, style)
}

// SetColWidth sets the width for a single column specified by colRef, eg "A", "BA" etc
func (f *File) SetColWidth(colRef string, width int) {
	f.XLSX.SetColWidth(f.SheetName, colRef, colRef, float64(width))
}

// SetCellStyle applies a style to the specified cell grid
func (f *File) SetCellStyle(startCell, endCell, style string) {
	st, _ := f.XLSX.NewStyle(style)
	f.XLSX.SetCellStyle(f.SheetName, startCell, endCell, st)
}

// AddRow adds a row of data to the sheet
func (f *File) AddRow(data []interface{}) error {

	if len(data) != len(f.Columns) {
		return errors.New("number of data items does not equal the number of columns")
	}

	f.NextRow++
	for i, c := range f.Columns {
		cell := c.Ref + strconv.Itoa(f.NextRow) // eg "A1", "A2"... "AA26"
		f.XLSX.SetCellValue(f.SheetName, cell, data[i])
	}

	return nil
}

// columnByName fetches a column by heading
func (f *File) colByHeading(heading string) (column, error) {
	var c column
	for _, c := range f.Columns {
		if c.Heading == heading {
			return c, nil
		}
	}
	return c, errors.New("column heading not found")
}

// AddError adds an error message to a separate sheet
func (f *File) AddError(id int, message string) {
	f.errorSheet() // Ensure error sheet
	row := []interface{}{id, message}
	axis := "A" + strconv.Itoa(f.NextErrRow)
	f.XLSX.SetSheetRow(errorSheetName, axis, &row)
	f.NextErrRow++
}

// errorSheet will initialise the error sheet if it is not already there
func (f *File) errorSheet() {
	i := f.XLSX.GetSheetIndex(errorSheetName)
	if i > 0 { // exists already
		return
	}
	f.XLSX.NewSheet(errorSheetName)
	f.XLSX.SetCellValue(errorSheetName, "A1", "Record ID")
	f.XLSX.SetCellValue(errorSheetName, "B1", "Error message")
}

// columnRefs generates the specified number of column references - eg "A", "B" ... "Z", "AA", "AB" etc.
func columnRefs(numCols int) []string {

	result := []string{}
	xa := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	for i := 0; i < numCols; i++ {

		var colName string
		var colPrefix string

		set := int(math.Floor(float64(i) / float64(26)))
		if set > 0 {
			colPrefix = xa[set-1]
		}
		colName = colPrefix + xa[i-(set*26)]
		result = append(result, colName)
	}

	return result
}

// rowrefs returns a []string containing the cell references for a row, eg "A10", "B10", "C10" etc
func rowRefs(columnKeys []string, rowNumber int) []string {
	var refs []string
	rowNum := strconv.Itoa(rowNumber)
	for _, c := range columnKeys {
		r := c + rowNum
		refs = append(refs, r)
	}
	return refs
}
