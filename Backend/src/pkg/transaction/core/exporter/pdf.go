package exporter

import (
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	"github.com/phpdave11/gofpdf"
)


func CreatePDFReport(title string, transactions []entity.Transaction) *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "A3", "") // custom wider than A3
	pdf.SetMargins(10, 15, 10)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, title)
	pdf.Ln(12)

	headers := []string{
		"ID", "Merchant ID", "User ID", "Phone", "Type", "Medium",
		"Reference", "Comment", "Reference No", "Description", "Token",
		"Verified", "TTL", "Has Challenge", "Webhook Received", "Test",
		"Amount", "Fee", "Admin Net", "VAT", "Merchant Net", "Total",
		"Currency", "Created At", "Updated At", "Confirmed At",
		"Callback URL", "Success URL", "Failed URL",
	}

	var dataRows [][]string
	for _, t := range transactions {
		row := []string{
			t.Id.String(),
			t.MerchantId.String(),
			t.UserId.String(),
			t.PhoneNumber,
			string(t.Type),
			string(t.Medium),
			t.Reference,
			t.Comment,
			t.ReferenceNumber,
			t.Description,
			t.Token,
			boolToStr(t.Verified),
			fmt.Sprintf("%d", t.TTL),
			boolToStr(t.HasChallenge),
			boolToStr(t.WebhookReceived),
			boolToStr(t.Test),
			fmt.Sprintf("%.2f", t.BaseAmount),
			fmt.Sprintf("%.2f", t.FeeAmount),
			fmt.Sprintf("%.2f", t.AdminNet),
			fmt.Sprintf("%.2f", t.VatAmount),
			fmt.Sprintf("%.2f", t.MerchantNet),
			fmt.Sprintf("%.2f", t.TotalAmount),
			t.Currency,
			formatTime(t.CreatedAt),
			formatTime(t.UpdatedAt),
			formatTime(t.Confirm_Timestamp),
			t.CallbackURL,
			t.SuccessURL,
			t.FailedURL,
		}
		dataRows = append(dataRows, row)
	}

	pdf.SetFont("Arial", "", 7)
	pageWidth, _ := pdf.GetPageSize()
	marginLeft, _, _, _ := pdf.GetMargins()
	maxWidth := pageWidth - 2*marginLeft

	// Define min widths per column (in mm)
	minWidths := []float64{
		25, 25, 25, 30, 35, 35,
		35, 45, 45, 40, 40, 40,
		50, 50, 45, 40, 40,
		25, 35, 35, 35, 35, 35,
		45, 42, 42, 42,
		30, 30, 30,
	}

	// Calculate actual widths
	colWidths := make([]float64, len(headers))
	for col := range headers {
		maxColWidth := pdf.GetStringWidth(headers[col]) + 4
		for _, row := range dataRows {
			w := pdf.GetStringWidth(row[col]) + 4
			if w > maxColWidth {
				maxColWidth = w
			}
		}
		// Take max of calculated and min width
		if maxColWidth < minWidths[col] {
			colWidths[col] = minWidths[col]
		} else {
			colWidths[col] = maxColWidth
		}
	}

	// Scale down if total width exceeds page
	totalWidth := 0.0
	for _, w := range colWidths {
		totalWidth += w
	}
	if totalWidth > maxWidth {
		scale := maxWidth / totalWidth
		for i := range colWidths {
			colWidths[i] *= scale
		}
	}

	// Header
	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(220, 220, 220)
	pdf.SetDrawColor(180, 180, 180)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 6, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Body with multiline support
	pdf.SetFont("Arial", "", 6.8)
	lineHeight := 4.2
	for rowIndex, row := range dataRows {
		maxLines := 1
		for i, val := range row {
			lines := pdf.SplitLines([]byte(val), colWidths[i])
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
		}
		cellHeight := float64(maxLines) * lineHeight

		if rowIndex%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		x := pdf.GetX()
		y := pdf.GetY()
		for i, val := range row {
			pdf.Rect(x, y, colWidths[i], cellHeight, "D")
			pdf.MultiCell(colWidths[i], lineHeight, val, "", "", false)
			x += colWidths[i]
			pdf.SetXY(x, y)
		}
		pdf.Ln(cellHeight)
	}

	return pdf
}

func truncate(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	if maxLen <= 3 {
		return str[:maxLen]
	}
	return str[:maxLen-3] + "..."
}

func boolToStr(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04")
}
