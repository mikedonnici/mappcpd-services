package activity

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"log"
)

// standard widths, heights, font sizes - for convenience
const (

	standardTextSize   = 12
	standardLineHeight = 6

	sectionHeadingTextSize   = 16
	sectionHeadingLineHeight = 12

	reportDetailTextSize   = 10
	reportDetailLineHeight = 4

	// addContext field name cells
	w1 = 30
	w2 = 140 // report main left cell
)

func PDFReport(reportData MemberActivityReport) {

	pdf := gofpdf.New("P", "mm", "A4", "")

	title := "CPD Activity Report..."
	pdf.SetTitle(title, false)
	pdf.SetAuthor("MappCPD PDF Generator", false)

	pdf.SetHeaderFunc(headerFunc(pdf, title))
	pdf.SetFooterFunc(footerFunc(pdf))
	pdf.AddPage()

	heading := fmt.Sprintf("%s (%s - %s)", reportData.Name, niceDate(reportData.StartDate), niceDate(reportData.EndDate))
	addSectionHeading(pdf, heading)
	addContext(pdf, reportData)

	addSectionHeading(pdf, "Summary")
	addSummary(pdf, reportData)

	addSectionHeading(pdf, "Detail")
	addDetail(pdf, reportData)

	err := pdf.OutputFileAndClose("report.pdf")
	if err != nil {
		fmt.Println(err)
	}
}

func addSectionHeading(pdf *gofpdf.Fpdf, subTitle string) {
	pdf.Ln(sectionHeadingLineHeight)
	pdf.SetFont("Arial", "B", sectionHeadingTextSize)
	pdf.MultiCell(0, sectionHeadingLineHeight, subTitle, "", "L", false)
}

func addContext(pdf *gofpdf.Fpdf, r MemberActivityReport) {
	pdf.SetFont("Arial", "", standardTextSize)
	pdf.Cell(w1, standardLineHeight, "Name:")
	pdf.Cell(w1, standardLineHeight, strconv.Itoa(r.MemberID))
	pdf.Ln(standardLineHeight)
	pdf.Cell(w1, standardLineHeight, "Member ID:")
	pdf.Cell(w1, standardLineHeight, strconv.Itoa(r.MemberID))
	pdf.Ln(standardLineHeight)
	pdf.Cell(w1, standardLineHeight, "Generated:")
	pdf.Cell(w1, standardLineHeight, time.Now().Format("02 Jan 2006 - 15:04 MST"))
	pdf.Ln(standardLineHeight)
}

func addSummary(pdf *gofpdf.Fpdf, r MemberActivityReport) {
	pdf.SetFont("Arial", "", standardTextSize)
	var total float64
	for _, a := range r.Activities {
		pdf.Cell(w2, standardLineHeight, a.ActivityName)
		pdf.CellFormat(w1, standardLineHeight, floatToString(a.CreditAwarded), "", 0, "R", false, 0, "")
		pdf.Ln(standardLineHeight)
		total += a.CreditAwarded
	}
	pdf.SetFont("Arial", "B", standardTextSize)
	pdf.Cell(w2, standardLineHeight, "Total")
	pdf.CellFormat(w1, standardLineHeight, floatToString(total), "", 0, "R", false, 0, "")
	pdf.Ln(standardLineHeight)
	pdf.Cell(w2, standardLineHeight, "Required")
	pdf.CellFormat(w1, standardLineHeight, floatToString(float64(r.CreditRequired)), "", 0, "R", false, 0, "")
	pdf.Ln(standardLineHeight)
}

func addDetail(pdf *gofpdf.Fpdf, r MemberActivityReport) {

	colWidths := []float64{22, 0, 22, 16, 16}
	colWidths[1] = pageDisplayWidth(pdf) - (colWidths[0] + colWidths[2] + colWidths[3] + colWidths[4])

	var total float64
	for _, a := range r.Activities {
		total += a.CreditAwarded
		addActivityDetailHeading(pdf, a)
		addActivityDetailColumnHeadings(pdf, colWidths)
		addActivityDetailRows(pdf, colWidths, a.Records)

		// addDetailFooter - total and capped
	}

	//pdf.SetFont("Arial", "", reportDetailTextSize)
	//var total float64
	//for _, a := range r.Activities {
	// activity summary
	//pdf.Cell(w2, standardLineHeight, a.ActivityName)
	//pdf.CellFormat(w1, standardLineHeight, floatToString(a.CreditAwarded), "1", 0, "R", false, 0, "")
	//pdf.Ln(standardLineHeight)
	//total += a.CreditAwarded

	//for _, r := range a.Records {
	//
	//	// calculate the height of the cell based on the length of column 2
	//	var rowHeight float64
	//	var c2LineCount int
	//	var c2Content string
	//
	//	// Type will sit on top of description line a sub-heading.
	//	if r.Type != "" {
	//		c2Content += r.Type + "\n"
	//		c2LineCount++
	//	}
	//
	//	descriptionLines := pdf.SplitLines([]byte(r.Description), widthC2)
	//	c2LineCount += len(descriptionLines)
	//	rowHeight = float64(c2LineCount) * standardLineHeight
	//
	//	for _, l := range descriptionLines {
	//		c2Content += string(l) + "\n"
	//	}
	//
	//	pdf.SetFont("Arial", "", reportDetailTextSize)
	//	pdf.CellFormat(widthC1, rowHeight, r.Date, "1", 0, "L", false, 0, "")
	//	//pdf.CellFormat(widthC2, rowHeight, c2Content, "1", 0, "L", false, 0, "")
	//	pdf.CellFormat(widthC3, rowHeight, "*", "1", 0, "C", false, 0, "")
	//	pdf.CellFormat(widthC4, rowHeight, floatToString(r.Quantity), "1", 0, "R", false, 0, "")
	//	pdf.CellFormat(widthC5, rowHeight, floatToString(r.Credit), "1", 0, "R", false, 0, "")
	//	pdf.Ln(8)
	//}

	//}
}

func pageDisplayWidth(pdf *gofpdf.Fpdf) float64 {
	pageWidth, _ := pdf.GetPageSize()
	pageMarginLeft, pageMarginRight, _, _ := pdf.GetMargins()
	pageDisplayWidth := pageWidth - (pageMarginLeft + pageMarginRight)

	return pageDisplayWidth
}

func addActivityDetailRows(pdf *gofpdf.Fpdf, colWidths []float64, rows []activity) {
	pdf.SetFont("Arial", "", reportDetailTextSize)

	for _, r := range rows {
		pdf.CellFormat(colWidths[0], reportDetailLineHeight, niceDate(r.Date), "0", 0, "L", false, 0, "")

		nextCellY := pdf.GetY() // go here are multicell
		if r.Type != "" {
			r.Description = r.Type + "\n" + r.Description
		}
		pdf.MultiCell(colWidths[1], reportDetailLineHeight, r.Description, "0", "L", false)
		nextRowY := pdf.GetY() // go here when row ends
		pdf.SetY(nextCellY)
		pdf.CellFormat(colWidths[0]+colWidths[1], reportDetailLineHeight, "", "0", 0, "C", false, 0, "")

		pdf.CellFormat(colWidths[2], reportDetailLineHeight, "?", "0", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], reportDetailLineHeight, floatToString(r.Quantity), "0", 0, "R", false, 0, "")
		pdf.CellFormat(colWidths[4], reportDetailLineHeight, floatToString(r.Credit), "0", 1, "R", false, 0, "")
		pdf.SetY(nextRowY)
		horizontalLine(pdf, pageDisplayWidth(pdf))
	}

	pdf.Ln(standardLineHeight)
}

func addActivityDetailHeading(pdf *gofpdf.Fpdf, a activityReport) {
	pdf.SetFont("Arial", "B", standardTextSize)
	pdf.MultiCell(0, standardLineHeight, a.ActivityName, "0", "L", false)
	pdf.SetFont("Arial", "", reportDetailTextSize)
	pdf.MultiCell(0, standardLineHeight, "Max credit: "+floatToString(a.MaxCredit), "0", "L", false)
}

func horizontalLine(pdf *gofpdf.Fpdf, width float64) {
	pdf.Ln(2)
	pdf.MultiCell(width, 2, "", "B", "C", false)
	pdf.Ln(2)
}

func addActivityDetailColumnHeadings(pdf *gofpdf.Fpdf, colWidths []float64) {
	pdf.SetFont("Arial", "B", reportDetailTextSize)
	pdf.CellFormat(colWidths[0], standardLineHeight, "Date", "1", 0, "L", false, 0, "")
	pdf.CellFormat(colWidths[1], standardLineHeight, "Type / Detail", "1", 0, "L", false, 0, "")
	pdf.CellFormat(colWidths[2], standardLineHeight, "Evidence", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidths[3], standardLineHeight, "Units", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidths[4], standardLineHeight, "Credit", "1", 1, "C", false, 0, "")
	pdf.Ln(reportDetailLineHeight/2)
}

func floatToString(n float64) string {
	return strconv.FormatFloat(n, 'f', 2, 64)
}

func niceDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error passing date -", err)
		return date
	}
	return t.Format("02 Jan 2006")
}

func headerFunc(pdf *gofpdf.Fpdf, titleStr string) func() {
	return func() {
		// Arial bold 15
		pdf.SetFont("Arial", "B", sectionHeadingTextSize)
		// Calculate width of title and position
		wd := pdf.GetStringWidth(titleStr) + 6
		pdf.SetX((210 - wd) / 2)
		// Colors of frame, background and text
		pdf.SetDrawColor(0, 80, 180)
		pdf.SetFillColor(230, 230, 0)
		// Thickness of frame (1 mm)
		pdf.SetLineWidth(1)
		// Title
		pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
		// Line break
		pdf.Ln(10)
	}
}

func footerFunc(pdf *gofpdf.Fpdf) func() {

	return func() {
		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 8)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	}
}
