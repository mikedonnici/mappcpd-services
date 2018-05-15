package cpd

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// standard widths, heights, font sizes - for convenience
const (
	text12  = 12
	height7 = 7

	text16   = 16
	height12 = 12

	text10  = 10
	height4 = 4

	width30  = 30
	width140 = 140
)

// PDFReport generates a PDF report and writes it to w
func PDFReport(reportData MemberActivityReport, w io.Writer) error {

	pdf := initPDF()
	addPageHeaderImage(pdf)
	addContextSection(pdf, reportData)
	addSummarySection(pdf, reportData)
	addDetailSection(pdf, reportData)

	return pdf.Output(w)
}

func initPDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	title := "CPD Activity Report..."
	pdf.SetTitle(title, false)
	pdf.SetAuthor("MappCPD PDF Generator", false)
	pdf.SetHeaderFunc(headerFunc(pdf))
	pdf.SetFooterFunc(footerFunc(pdf))
	pdf.SetDrawColor(221, 221, 221) // for borders
	pdf.AddPage()

	return pdf
}

func addContextSection(pdf *gofpdf.Fpdf, reportData MemberActivityReport) {
	heading := fmt.Sprintf("%s (%s - %s)", reportData.ReportName, niceDate(reportData.StartDate), niceDate(reportData.EndDate))
	addSectionHeading(pdf, heading)
	addContext(pdf, reportData)
}

func addSummarySection(pdf *gofpdf.Fpdf, reportData MemberActivityReport) {
	addSectionHeading(pdf, "Summary")
	addSummary(pdf, reportData)
}

func addDetailSection(pdf *gofpdf.Fpdf, reportData MemberActivityReport) {
	addSectionHeading(pdf, "Detail")
	addDetail(pdf, reportData)
}

func addSectionHeading(pdf *gofpdf.Fpdf, subTitle string) {
	pdf.Ln(height12)
	pdf.SetFont("Arial", "B", text16)
	pdf.MultiCell(0, height12, subTitle, "", "L", false)
}

func addContext(pdf *gofpdf.Fpdf, r MemberActivityReport) {
	pdf.SetFont("Arial", "", text12)
	pdf.Cell(width30, height7, "Name:")
	pdf.Cell(width30, height7, strconv.Itoa(r.MemberID))
	pdf.Ln(height7)
	pdf.Cell(width30, height7, "Member ID:")
	pdf.Cell(width30, height7, strconv.Itoa(r.MemberID))
	pdf.Ln(height7)
	pdf.Cell(width30, height7, "Generated:")
	pdf.Cell(width30, height7, time.Now().Format("02 Jan 2006 - 15:04 MST"))
	pdf.Ln(height7)
}

func addSummary(pdf *gofpdf.Fpdf, r MemberActivityReport) {
	pdf.SetFont("Arial", "", text12)
	var total float64
	for _, a := range r.Activities {
		pdf.Cell(width140, height7, a.ActivityName)
		pdf.CellFormat(width30, height7, floatToString(a.CreditAwarded), "", 0, "R", false, 0, "")
		pdf.Ln(height7)
		total += a.CreditAwarded
	}
	addRowDividerLine(pdf, 0)
	pdf.SetFont("Arial", "B", text12)
	pdf.CellFormat(width140, height7, "Total:", "", 0, "R", false, 0, "")
	pdf.CellFormat(width30, height7, floatToString(total), "", 0, "R", false, 0, "")
	pdf.CellFormat(width140, height7, "Required:", "", 0, "R", false, 0, "")
	pdf.CellFormat(width30, height7, floatToString(float64(r.CreditRequired)), "", 0, "R", false, 0, "")
	pdf.Ln(height7)
}

func addDetail(pdf *gofpdf.Fpdf, r MemberActivityReport) {

	colWidths := []float64{22, 0, 16, 16, 16}
	colWidths[1] = pageDisplayWidth(pdf) - (colWidths[0] + colWidths[2] + colWidths[3] + colWidths[4])

	var total float64
	for _, a := range r.Activities {
		total += a.CreditAwarded
		addActivityDetailHeading(pdf, a)
		addActivityDetailColumnHeadings(pdf, colWidths)
		addActivityDetailRows(pdf, colWidths, a.Records)
		// addDetailFooter - total and capped
	}
}

func pageDisplayWidth(pdf *gofpdf.Fpdf) float64 {
	pageWidth, _ := pdf.GetPageSize()
	pageMarginLeft, pageMarginRight, _, _ := pdf.GetMargins()
	pageDisplayWidth := pageWidth - (pageMarginLeft + pageMarginRight)

	return pageDisplayWidth
}

func addActivityDetailRows(pdf *gofpdf.Fpdf, colWidths []float64, records []activityRecord) {

	pdf.SetFont("Arial", "", text10)

	for _, r := range records {
		pdf.CellFormat(colWidths[0], height4, niceDate(r.Date), "0", 0, "L", false, 0, "")

		nextCellY := pdf.GetY() // go here after multicell
		if r.Type != "" {
			r.Description = r.Type + " : " + r.Description
		}
		pdf.MultiCell(colWidths[1], height4, r.Description, "0", "L", false)
		nextRowY := pdf.GetY() // go here when row ends
		pdf.SetY(nextCellY)
		pdf.CellFormat(colWidths[0]+colWidths[1], height4, "", "0", 0, "C", false, 0, "")

		pdf.CellFormat(colWidths[2], height4, "?", "0", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], height4, floatToString(r.Quantity), "0", 0, "R", false, 0, "")
		pdf.CellFormat(colWidths[4], height4, floatToString(r.Credit), "0", 1, "R", false, 0, "")
		pdf.SetY(nextRowY)
		addRowDividerLine(pdf, pageDisplayWidth(pdf))
	}

	pdf.Ln(height7)
}

func addRowDividerLine(pdf *gofpdf.Fpdf, width float64) {
	pdf.Ln(2)
	pdf.MultiCell(width, 2, "", "B", "C", false)
	pdf.Ln(4)
}

func addActivityDetailHeading(pdf *gofpdf.Fpdf, a activityReport) {
	pdf.SetFont("Arial", "B", text12)
	pdf.MultiCell(0, height7, a.ActivityName, "0", "L", false)
	pdf.SetFont("Arial", "", text10)
	pdf.MultiCell(0, height7, "Max credit: "+floatToString(a.MaxCredit), "0", "L", false)
}

func addActivityDetailColumnHeadings(pdf *gofpdf.Fpdf, colWidths []float64) {
	pdf.SetFont("Arial", "B", text10)
	pdf.CellFormat(colWidths[0], height7, "Date", "B", 0, "L", false, 0, "")
	pdf.CellFormat(colWidths[1], height7, "Type / Detail", "B", 0, "L", false, 0, "")
	pdf.CellFormat(colWidths[2], height7, "Evidence", "B", 0, "C", false, 0, "")
	pdf.CellFormat(colWidths[3], height7, "Units", "B", 0, "R", false, 0, "")
	pdf.CellFormat(colWidths[4], height7, "Credit", "B", 1, "R", false, 0, "")
	pdf.Ln(height4 / 2)
}

func addPageHeaderImage(pdf *gofpdf.Fpdf) {

	res, err := http.Get("https://d1cbfvxg6albaj.cloudfront.net/pdf/header.jpg")
	if err != nil {
		return
	}
	defer res.Body.Close()

	pdf.RegisterImageReader("header.jpg", "JPG", res.Body)
	pdf.Image("header.jpg", 0, 0, 210, 0, false, "", 0, "")
	pdf.Ln(height7 * 2)
}

func headerFunc(pdf *gofpdf.Fpdf) func() {

	return func() {
		text := "CPD Report generated " + time.Now().Format("02 Jan 2006 - 15:04 MST")
		pdf.SetY(5)
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(30, 10, text, "0", 0, "L", false, 0, "")
		pdf.Ln(height7 * 2)
	}
}

func footerFunc(pdf *gofpdf.Fpdf) func() {

	return func() {
		text := fmt.Sprintf("Page %d", pdf.PageNo())
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(0, 10, text, "", 0, "R", false, 0, "")
	}
}

func floatToString(n float64) string {
	return strconv.FormatFloat(n, 'f', 2, 64)
}

// niceDate returns unmodified date string on error
func niceDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.Format("02 Jan 06")
}
