package utils

import (
	"archive/zip"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	"github.com/xuri/excelize/v2"
)

// OptionalBool format bool instance with nil value
func OptionalBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

// OptionalString formats string instance with nil value
func OptionalString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// OptionalTime formats time instance with nil value
func OptionalTime(t *time.Time) string {
	if t != nil {
		return t.Format(time.RFC3339)
	}
	return ""
}

// optionalUUID formats uuid instance with nil value
func OptionalUUID(u *uuid.UUID) string {
	if u != nil {
		return u.String()
	}
	return ""
}

// ToNullUUID converts uuid instance to uuid.NullUUID
func ToNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u != nil {
		return uuid.NullUUID{UUID: *u, Valid: true}
	}
	return uuid.NullUUID{Valid: false}
}

// ToNullString converts string instance to sql.NullString
func ToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

// ToNullBool converts bool instance to sql.NullBool
func ToNullBool(b *bool) sql.NullBool {
	if b != nil {
		return sql.NullBool{Bool: *b, Valid: true}
	}
	return sql.NullBool{Valid: false}
}

// ToNullTime converts time instance to sql.NullTime
func ToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

// writeCSVToZip appends CSV file to a single zip file
func writeCSVToZip(zipWriter *zip.Writer, filename string, writeFunc func(w *csv.Writer) error) error {
	f, err := zipWriter.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s in zip: %w", filename, err)
	}
	writer := csv.NewWriter(f)
	if err := writeFunc(writer); err != nil {
		return fmt.Errorf("failed writing %s: %w", filename, err)
	}
	writer.Flush()
	return writer.Error()
}

// FilterMerchantDetailsJSON filters field from MerchantDetails instance
func FilterMerchantDetailsJSON(data entity.MerchantDetails, includeFields []string) (entity.MerchantDetails, error) {
	jsonData, _ := json.Marshal(data)
	var asMap map[string]interface{}
	json.Unmarshal(jsonData, &asMap)

	// Filter fields
	filteredMap := map[string]interface{}{}
	fieldSet := make(map[string]struct{})
	for _, f := range includeFields {
		fieldSet[f] = struct{}{}
	}

	for k, v := range asMap {
		if _, ok := fieldSet[k]; ok {
			filteredMap[k] = v
		}
	}

	// Convert back
	filteredJSON, _ := json.Marshal(filteredMap)
	var result entity.MerchantDetails
	err := json.Unmarshal(filteredJSON, &result)
	return result, err
}

// WriteMerchantDetailsToCSVZip writes MerchantDetails into separate CSVs and zips them
func WriteMerchantDetailsToCSVZip(data []entity.MerchantDetails, zipFilename string) error {
	publicDir := "public"
	zipPath := filepath.Join(publicDir, zipFilename)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Helper to check if we should skip writing a section
	shouldWrite := func(count int) bool {
		return count > 0
	}

	// ------------------------
	// Write Merchants.csv
	// ------------------------
	if err := writeCSVToZip(zipWriter, "Merchants.csv", func(w *csv.Writer) error {
		header := []string{
			"ID", "UserID", "LegalName", "TradingName", "BusinessRegistrationNumber",
			"TaxIdentificationNumber", "BusinessType", "IndustryCategory", "IsBettingCompany",
			"LotteryCertificateNumber", "WebsiteURL", "EstablishedDate", "CreatedAt", "UpdatedAt", "Status",
		}
		if err := w.Write(header); err != nil {
			return err
		}
		for _, md := range data {
			m := md.Merchant
			row := []string{
				m.ID.String(),
				m.UserID.String(),
				m.LegalName,
				OptionalString(m.TradingName),
				m.BusinessRegistrationNumber,
				m.TaxIdentificationNumber,
				m.BusinessType,
				OptionalString(m.IndustryCategory),
				fmt.Sprintf("%v", m.IsBettingCompany),
				OptionalString(m.LotteryCertificateNumber),
				OptionalString(m.WebsiteURL),
				OptionalTime(m.EstablishedDate),
				m.CreatedAt.Format(time.RFC3339),
				m.UpdatedAt.Format(time.RFC3339),
				string(m.Status),
			}
			if err := w.Write(row); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// ------------------------
	// Write Addresses.csv
	// ------------------------
	addressCount := 0
	for _, md := range data {
		addressCount += len(md.Addresses)
	}
	if shouldWrite(addressCount) {
		if err := writeCSVToZip(zipWriter, "Addresses.csv", func(w *csv.Writer) error {
			header := []string{
				"ID", "MerchantID", "AddressType", "StreetAddress1", "StreetAddress2",
				"City", "Region", "PostalCode", "Country", "IsPrimary", "CreatedAt", "UpdatedAt",
			}
			if err := w.Write(header); err != nil {
				return err
			}
			for _, md := range data {
				for _, addr := range md.Addresses {
					row := []string{
						addr.ID.String(),
						addr.MerchantID.String(),
						addr.AddressType,
						addr.StreetAddress1,
						OptionalString(addr.StreetAddress2),
						addr.City,
						addr.Region,
						OptionalString(addr.PostalCode),
						addr.Country,
						fmt.Sprintf("%v", addr.IsPrimary),
						addr.CreatedAt.Format(time.RFC3339),
						addr.UpdatedAt.Format(time.RFC3339),
					}
					if err := w.Write(row); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// ------------------------
	// Write Contacts.csv
	// ------------------------
	contactCount := 0
	for _, md := range data {
		contactCount += len(md.Contacts)
	}
	if shouldWrite(contactCount) {
		if err := writeCSVToZip(zipWriter, "Contacts.csv", func(w *csv.Writer) error {
			header := []string{
				"ID", "MerchantID", "ContactType", "FirstName", "LastName", "Email",
				"PhoneNumber", "Position", "IsVerified", "CreatedAt", "UpdatedAt",
			}
			if err := w.Write(header); err != nil {
				return err
			}
			for _, md := range data {
				for _, c := range md.Contacts {
					row := []string{
						c.ID.String(),
						c.MerchantID.String(),
						c.ContactType,
						c.FirstName,
						c.LastName,
						c.Email,
						c.PhoneNumber,
						OptionalString(c.Position),
						fmt.Sprintf("%v", c.IsVerified),
						c.CreatedAt.Format(time.RFC3339),
						c.UpdatedAt.Format(time.RFC3339),
					}
					if err := w.Write(row); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// ------------------------
	// Write Documents.csv
	// ------------------------
	documentCount := 0
	for _, md := range data {
		documentCount += len(md.Documents)
	}
	if shouldWrite(documentCount) {
		if err := writeCSVToZip(zipWriter, "Documents.csv", func(w *csv.Writer) error {
			header := []string{
				"ID", "MerchantID", "DocumentType", "DocumentNumber", "FileURL", "FileHash",
				"VerifiedBy", "VerifiedAt", "Status", "RejectionReason", "CreatedAt", "UpdatedAt",
			}
			if err := w.Write(header); err != nil {
				return err
			}
			for _, md := range data {
				for _, doc := range md.Documents {
					row := []string{
						doc.ID.String(),
						doc.MerchantID.String(),
						doc.DocumentType,
						OptionalString(doc.DocumentNumber),
						doc.FileURL,
						OptionalString(doc.FileHash),
						OptionalUUID(doc.VerifiedBy),
						OptionalTime(doc.VerifiedAt),
						doc.Status,
						OptionalString(doc.RejectionReason),
						doc.CreatedAt.Format(time.RFC3339),
						doc.UpdatedAt.Format(time.RFC3339),
					}
					if err := w.Write(row); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// ------------------------
	// Write BankAccounts.csv
	// ------------------------
	bankAccountCount := 0
	for _, md := range data {
		bankAccountCount += len(md.BankAccounts)
	}
	if shouldWrite(bankAccountCount) {
		if err := writeCSVToZip(zipWriter, "BankAccounts.csv", func(w *csv.Writer) error {
			header := []string{
				"ID", "MerchantID", "AccountHolderName", "BankName", "BankCode", "BranchCode",
				"AccountNumber", "AccountType", "Currency", "IsPrimary", "IsVerified",
				"VerificationDocumentID", "CreatedAt", "UpdatedAt",
			}
			if err := w.Write(header); err != nil {
				return err
			}
			for _, md := range data {
				for _, ba := range md.BankAccounts {
					row := []string{
						ba.ID.String(),
						ba.MerchantID.String(),
						ba.AccountHolderName,
						ba.BankName,
						ba.BankCode,
						OptionalString(ba.BranchCode),
						ba.AccountNumber,
						ba.AccountType,
						ba.Currency,
						fmt.Sprintf("%v", ba.IsPrimary),
						fmt.Sprintf("%v", ba.IsVerified),
						OptionalUUID(ba.VerificationDocumentID),
						ba.CreatedAt.Format(time.RFC3339),
						ba.UpdatedAt.Format(time.RFC3339),
					}
					if err := w.Write(row); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// ------------------------
	// Write Settings.csv
	// ------------------------
	settingsCount := 0
	for _, md := range data {
		if md.Settings != nil {
			settingsCount++
		}
	}
	if shouldWrite(settingsCount) {
		if err := writeCSVToZip(zipWriter, "Settings.csv", func(w *csv.Writer) error {
			header := []string{
				"MerchantID", "DefaultCurrency", "DefaultLanguage", "CheckoutTheme",
				"EnableWebhooks", "WebhookURL", "WebhookSecret", "AutoSettlement",
				"SettlementFrequency", "RiskSettings", "CreatedAt", "UpdatedAt",
			}
			if err := w.Write(header); err != nil {
				return err
			}
			for _, md := range data {
				if md.Settings != nil {
					s := md.Settings
					row := []string{
						s.MerchantID.String(),
						s.DefaultCurrency,
						s.DefaultLanguage,
						OptionalString(s.CheckoutTheme),
						fmt.Sprintf("%v", s.EnableWebhooks),
						OptionalString(s.WebhookURL),
						OptionalString(s.WebhookSecret),
						fmt.Sprintf("%v", s.AutoSettlement),
						s.SettlementFrequency,
						OptionalString(s.RiskSettings),
						s.CreatedAt.Format(time.RFC3339),
						s.UpdatedAt.Format(time.RFC3339),
					}
					if err := w.Write(row); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	fmt.Printf("ZIP archive created at: %s\n", zipPath)
	return nil
}

// WriteMerchantDetailsToJSON writes an array of MerchantDetails to a JSON file
func WriteMerchantDetailsToJSON(data []entity.MerchantDetails, filename string) error {
	publicDir := "public"
	jsonPath := filepath.Join(publicDir, filename)
	file, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	fmt.Printf("JSON file created at: %s\n", jsonPath)
	return nil
}

// WriteMerchantDetailsToXLSX writes MerchantDetails into an Excel file with separate sheets
func WriteMerchantDetailsToXLSX(data []entity.MerchantDetails, filename string) error {
	publicDir := "public"
	xlsxPath := filepath.Join(publicDir, filename)
	f := excelize.NewFile()

	err := f.DeleteSheet("Sheet1")
	if err != nil {
		return fmt.Errorf("failed to delete default sheet: %w", err)
	}

	// ------------------------
	// Sheet: Merchants
	// ------------------------
	merchantSheet := "Merchants"
	f.NewSheet(merchantSheet)
	merchantHeaders := []string{
		"ID", "UserID", "LegalName", "TradingName", "BusinessRegistrationNumber",
		"TaxIdentificationNumber", "BusinessType", "IndustryCategory", "IsBettingCompany",
		"LotteryCertificateNumber", "WebsiteURL", "EstablishedDate", "CreatedAt", "UpdatedAt", "Status",
	}
	for col, header := range merchantHeaders {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(merchantSheet, cell, header)
	}
	for rowIdx, md := range data {
		row := rowIdx + 2
		m := md.Merchant
		values := []interface{}{
			m.ID.String(),
			m.UserID.String(),
			m.LegalName,
			OptionalString(m.TradingName),
			m.BusinessRegistrationNumber,
			m.TaxIdentificationNumber,
			m.BusinessType,
			OptionalString(m.IndustryCategory),
			m.IsBettingCompany,
			OptionalString(m.LotteryCertificateNumber),
			OptionalString(m.WebsiteURL),
			OptionalTime(m.EstablishedDate),
			m.CreatedAt.Format(time.RFC3339),
			m.UpdatedAt.Format(time.RFC3339),
			string(m.Status),
		}
		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row)
			f.SetCellValue(merchantSheet, cell, val)
		}
	}
	activeSheet, err := f.GetSheetIndex(merchantSheet)

	if err != nil {
		return fmt.Errorf("failed to get active sheet: %w", err)
	}

	f.SetActiveSheet(activeSheet)

	// ------------------------
	// Sheet: Addresses
	// ------------------------
	addressCount := 0
	for _, md := range data {
		addressCount += len(md.Addresses)
	}
	if addressCount > 0 {
		addressSheet := "Addresses"
		f.NewSheet(addressSheet)
		addressHeaders := []string{
			"ID", "MerchantID", "AddressType", "StreetAddress1", "StreetAddress2",
			"City", "Region", "PostalCode", "Country", "IsPrimary", "CreatedAt", "UpdatedAt",
		}
		for col, header := range addressHeaders {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(addressSheet, cell, header)
		}
		row := 2
		for _, md := range data {
			for _, addr := range md.Addresses {
				values := []interface{}{
					addr.ID.String(),
					addr.MerchantID.String(),
					addr.AddressType,
					addr.StreetAddress1,
					OptionalString(addr.StreetAddress2),
					addr.City,
					addr.Region,
					OptionalString(addr.PostalCode),
					addr.Country,
					addr.IsPrimary,
					addr.CreatedAt.Format(time.RFC3339),
					addr.UpdatedAt.Format(time.RFC3339),
				}
				for col, val := range values {
					cell, _ := excelize.CoordinatesToCellName(col+1, row)
					f.SetCellValue(addressSheet, cell, val)
				}
				row++
			}
		}
	}

	// ------------------------
	// Sheet: Contacts
	// ------------------------
	contactCount := 0
	for _, md := range data {
		contactCount += len(md.Contacts)
	}
	if contactCount > 0 {
		contactSheet := "Contacts"
		f.NewSheet(contactSheet)
		contactHeaders := []string{
			"ID", "MerchantID", "ContactType", "FirstName", "LastName", "Email",
			"PhoneNumber", "Position", "IsVerified", "CreatedAt", "UpdatedAt",
		}
		for col, header := range contactHeaders {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(contactSheet, cell, header)
		}
		row := 2
		for _, md := range data {
			for _, c := range md.Contacts {
				values := []interface{}{
					c.ID.String(),
					c.MerchantID.String(),
					c.ContactType,
					c.FirstName,
					c.LastName,
					c.Email,
					c.PhoneNumber,
					OptionalString(c.Position),
					c.IsVerified,
					c.CreatedAt.Format(time.RFC3339),
					c.UpdatedAt.Format(time.RFC3339),
				}
				for col, val := range values {
					cell, _ := excelize.CoordinatesToCellName(col+1, row)
					f.SetCellValue(contactSheet, cell, val)
				}
				row++
			}
		}
	}

	// ------------------------
	// Sheet: Documents
	// ------------------------
	documentCount := 0
	for _, md := range data {
		documentCount += len(md.Documents)
	}
	if documentCount > 0 {
		docSheet := "Documents"
		f.NewSheet(docSheet)
		docHeaders := []string{
			"ID", "MerchantID", "DocumentType", "DocumentNumber", "FileURL", "FileHash",
			"VerifiedBy", "VerifiedAt", "Status", "RejectionReason", "CreatedAt", "UpdatedAt",
		}
		for col, header := range docHeaders {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(docSheet, cell, header)
		}
		row := 2
		for _, md := range data {
			for _, doc := range md.Documents {
				values := []interface{}{
					doc.ID.String(),
					doc.MerchantID.String(),
					doc.DocumentType,
					OptionalString(doc.DocumentNumber),
					doc.FileURL,
					OptionalString(doc.FileHash),
					OptionalUUID(doc.VerifiedBy),
					OptionalTime(doc.VerifiedAt),
					doc.Status,
					OptionalString(doc.RejectionReason),
					doc.CreatedAt.Format(time.RFC3339),
					doc.UpdatedAt.Format(time.RFC3339),
				}
				for col, val := range values {
					cell, _ := excelize.CoordinatesToCellName(col+1, row)
					f.SetCellValue(docSheet, cell, val)
				}
				row++
			}
		}
	}

	// ------------------------
	// Sheet: BankAccounts
	// ------------------------
	bankCount := 0
	for _, md := range data {
		bankCount += len(md.BankAccounts)
	}
	if bankCount > 0 {
		bankSheet := "BankAccounts"
		f.NewSheet(bankSheet)
		bankHeaders := []string{
			"ID", "MerchantID", "AccountHolderName", "BankName", "BankCode", "BranchCode",
			"AccountNumber", "AccountType", "Currency", "IsPrimary", "IsVerified",
			"VerificationDocumentID", "CreatedAt", "UpdatedAt",
		}
		for col, header := range bankHeaders {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(bankSheet, cell, header)
		}
		row := 2
		for _, md := range data {
			for _, ba := range md.BankAccounts {
				values := []interface{}{
					ba.ID.String(),
					ba.MerchantID.String(),
					ba.AccountHolderName,
					ba.BankName,
					ba.BankCode,
					OptionalString(ba.BranchCode),
					ba.AccountNumber,
					ba.AccountType,
					ba.Currency,
					ba.IsPrimary,
					ba.IsVerified,
					OptionalUUID(ba.VerificationDocumentID),
					ba.CreatedAt.Format(time.RFC3339),
					ba.UpdatedAt.Format(time.RFC3339),
				}
				for col, val := range values {
					cell, _ := excelize.CoordinatesToCellName(col+1, row)
					f.SetCellValue(bankSheet, cell, val)
				}
				row++
			}
		}
	}

	// ------------------------
	// Sheet: Settings
	// ------------------------
	settingsCount := 0
	for _, md := range data {
		if md.Settings != nil {
			settingsCount++
		}
	}
	if settingsCount > 0 {
		settingsSheet := "Settings"
		f.NewSheet(settingsSheet)
		settingsHeaders := []string{
			"MerchantID", "DefaultCurrency", "DefaultLanguage", "CheckoutTheme",
			"EnableWebhooks", "WebhookURL", "WebhookSecret", "AutoSettlement",
			"SettlementFrequency", "RiskSettings", "CreatedAt", "UpdatedAt",
		}
		for col, header := range settingsHeaders {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(settingsSheet, cell, header)
		}
		row := 2
		for _, md := range data {
			if md.Settings != nil {
				s := md.Settings
				values := []interface{}{
					s.MerchantID.String(),
					s.DefaultCurrency,
					s.DefaultLanguage,
					OptionalString(s.CheckoutTheme),
					s.EnableWebhooks,
					OptionalString(s.WebhookURL),
					OptionalString(s.WebhookSecret),
					s.AutoSettlement,
					s.SettlementFrequency,
					OptionalString(s.RiskSettings),
					s.CreatedAt.Format(time.RFC3339),
					s.UpdatedAt.Format(time.RFC3339),
				}
				for col, val := range values {
					cell, _ := excelize.CoordinatesToCellName(col+1, row)
					f.SetCellValue(settingsSheet, cell, val)
				}
				row++
			}
		}
	}

	// Save Excel file
	if err := f.SaveAs(xlsxPath); err != nil {
		return fmt.Errorf("failed to save XLSX: %w", err)
	}

	fmt.Printf("Excel file created at: %s\n", xlsxPath)
	return nil
}
