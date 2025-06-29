package usecase

import (
	"bytes"
	"context"
	"crypto/tls"
	stdErrors "errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/socialpay/socialpay/src/pkg/merchants/core/entity"
	"github.com/socialpay/socialpay/src/pkg/merchants/errors"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	"github.com/lib/pq"

	"github.com/google/uuid"
	//import procedeure "github.com/socialpay/socialpay/src/pkg/key/adapter/controller/procedure"
)

func (u *Usecase) CreateMerchant(merchant entity.Merchant) (*entity.Merchant, error) {

	if err := merchant.Validate(); err != nil {

		u.log.Println("INFO::validation err :", err)

		return nil, err

	}
	// Save the key pair to the repository
	err := u.repo.Save(&merchant)
	if err != nil {
		u.log.Printf("DB_ERR::ERR WHILE SAVING MERCHANT ::%v", err.Error())

		var pqErr *pq.Error

		if stdErrors.As(err, &pqErr) {

			switch pqErr.Code {
			case "23505": // unique_violation
				return nil, errors.Error{
					Type:    "CONFLICT_ERROR",
					Message: "Merchant already exists ",
					Code:    http.StatusConflict,
				}
			case "23503": // foreign_key_violation
				return nil, errors.Error{
					Type:    "INVALID_REFERENCE",
					Message: "Invalid foreign key reference",
					Code:    http.StatusBadRequest,
				}
			}

		}
		err = errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: "Failed to save merchant please try again",
			Code:    http.StatusInternalServerError,
		}
		return nil, err
	}

	return &merchant, nil
}

func (u *Usecase) GetMerchantByUserID(userID uuid.UUID) (*entity.Merchant, error) {
	// Find the key pair by user ID from the repository
	merchant, err := u.repo.FindByUserID(userID)
	if err != nil {
		u.log.Printf("ERR WHILE FINDING MERCHANT BY USER ID::%v", err.Error())
		err = errors.Error{
			Type:    "NOT_FOUND",
			Message: "Merchant not found",
			Code:    http.StatusNotFound,
		}
		return nil, err
	}

	return merchant, nil
}

func (u *Usecase) GetMerchantDetails(userID uuid.UUID) (*entity.MerchantDetails, error) {
	merchant, err := u.repo.GetMerchantDetails(userID)
	if err != nil {
		u.log.Printf("Error getting merchant details: %v", err)
		return nil, errors.Error{
			Type:    "DATABASE_ERROR",
			Message: "Failed to retrieve merchant details",
			Code:    http.StatusInternalServerError,
		}
	}
	return merchant, nil
}

func (u *Usecase) GetMerchants() ([]entity.Merchant, error) {
	merchants, err := u.repo.GetMerchants()
	if err != nil {
		u.log.Printf("Error while fetching merchants: %v", err)
		return nil, errors.Error{
			Type:    "DATABASE_ERROR",
			Message: "Failed to retrieve merchants",
			Code:    http.StatusInternalServerError,
		}
	}

	if len(merchants) == 0 {
		u.log.Println("No merchants found")
		return nil, errors.Error{
			Type:    "NOT_FOUND",
			Message: "No merchants available",
			Code:    http.StatusNotFound,
		}
	}

	return merchants, nil
}

// usecase/merchant.go
func (u *Usecase) UpdateMerchantStatus(merchantID uuid.UUID, status entity.MerchantStatus) error {

	// Retrieve the merchant details for logging
	merchant, err := u.repo.GetMerchantByID(merchantID)
	if err != nil {
		u.log.Printf("Error retrieving merchant details: %v", err)
		return errors.Error{
			Type:    "DATABASE_ERROR",
			Message: "Failed to retrieve merchant details",
			Code:    http.StatusInternalServerError,
		}
	}

	// Update status
	err = u.repo.UpdateMerchantStatus(merchantID, status)
	if err != nil {
		u.log.Printf("Error updating merchant status: %v", err)
		return errors.Error{
			Type:    "DATABASE_ERROR",
			Message: "Failed to update merchant status",
			Code:    http.StatusInternalServerError,
		}
	}

	// Log status change
	u.log.Printf("Merchant %s status changed from %s to %s", merchantID, merchant.Status, status)

	return nil
}

func (u *Usecase) DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error {
	// Validate merchant exists first
	merchant, err := u.repo.GetMerchantByID(merchantID)
	if err != nil {
		u.log.Printf("Error finding merchant to delete: %v", err)
		return errors.Error{
			Type:    "NOT_FOUND",
			Message: "Merchant not found",
			Code:    http.StatusNotFound,
		}
	}
	// prevent deletion of active merchants:
	if merchant.Status == string(entity.StatusActive) {
		return errors.Error{
			Type:    "BUSINESS_RULE",
			Message: "Cannot delete active merchants",
			Code:    http.StatusBadRequest,
		}
	}

	// Perform the deletion
	if err := u.repo.DeleteMerchant(ctx, merchantID); err != nil {
		u.log.Printf("Error deleting merchant: %v", err)
		return errors.Error{
			Type:    "DATABASE_ERROR",
			Message: "Failed to delete merchant",
			Code:    http.StatusInternalServerError,
		}
	}

	u.log.Printf("Successfully deleted merchant %s", merchantID)
	return nil
}

func (u *Usecase) UpdateFullMerchant(ctx context.Context, merchantID uuid.UUID,
	merchant *entity.Merchant, address *entity.MerchantAdditionalInfo,
	documents []entity.MerchantDocument) error {

	// Validate merchant exists
	if _, err := u.repo.GetMerchantByID(merchantID); err != nil {
		u.log.Printf("Merchant not found: %v", err)
		return errors.Error{
			Type:    "NOT_FOUND",
			Message: "Merchant not found",
			Code:    http.StatusNotFound,
		}
	}

	// Business validation could go here
	if merchant.IsBettingClient && merchant.LoyaltyCertificate == "" {
		return errors.Error{
			Type:    "VALIDATION",
			Message: "Betting companies must provide lottery certificate",
			Code:    http.StatusBadRequest,
		}
	}

	// Perform update
	if err := u.repo.UpdateFullMerchant(ctx, merchantID, merchant, address, documents); err != nil {
		u.log.Printf("Failed to update merchant: %v", err)
		return errors.Error{
			Type:    "DATABASE_ERROR",
			Message: "Failed to update merchant",
			Code:    http.StatusInternalServerError,
		}
	}

	u.log.Printf("Successfully updated merchant %s", merchantID)
	return nil
}

func (u *Usecase) AddMerchantInfo(ctx context.Context, userId uuid.UUID,
	info entity.MerchantAdditionalInfo) error {

	// Validate the input
	if err := entity.ValidateMerchantAdditionalInfo(info); err != nil {
		u.log.Printf("ERR WHILE VALIDATING MERCHANT ADDITIONAL INFO::%v", err.Error())
		return err
	}
	// Find the key pair by user ID from the repository
	merchant, err := u.repo.FindByUserID(userId)

	if err != nil {
		u.log.Printf("ERR WHILE FINDING MERCHANT BY USER ID::%v", err.Error())

		err = errors.Error{
			Type:    "NOT_FOUND",
			Message: "Merchant not found",
			Code:    http.StatusNotFound,
		}

		return err
	}
	info.MerchantID, err = uuid.Parse(merchant.MerchantID)
	if err != nil {
		u.log.Printf("ERR WHILE PARSING MERCHANT ID::%v", err.Error())
		err = errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: "Merchant ID parsing error",
			Code:    http.StatusInternalServerError,
		}
		return err
	}
	// Save the key pair to the repository
	if err := u.repo.CreateMerchantAdditionalInfo(ctx, info); err != nil {
		u.log.Printf("ERR WHILE SAVING MERCHANT ADDITIONAL INFO::%v", err.Error())
		err = errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: "Failed to save merchant additional info please try again",
			Code:    http.StatusInternalServerError,
		}
		return err
	}
	return nil

}

func (u *Usecase) AddDocument(ctx context.Context, file multipart.File,
	fileHeader multipart.FileHeader, doc entity.MerchantDocument) error {

	var buf bytes.Buffer

	cld, err := InitCloudinary()

	if err != nil {
		return err
	}

	// buffer the file
	_, err = io.Copy(&buf, file)
	if err != nil {
		u.log.Printf("ERR WHILE BUFFERING FILE::%v", err.Error())
		err = errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: "Failed to buffer file",
			Code:    http.StatusInternalServerError,
		}
		return err
	}
	// Upload the file to Cloudinary
	res, err := cld.Upload.Upload(ctx, &buf, uploader.UploadParams{
		PublicID: fileHeader.Filename + uuid.NewString(),
		Folder:   "Doc",
	})

	if err != nil {
		u.log.Printf("ERR WHILE UPLOADING FILE TO CLOUDINARY::%v", err.Error())
		err = errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: "Failed to upload file to cloudinary ::" + err.Error(),
			Code:    http.StatusInternalServerError,
		}
		return err
	}

	doc.FileURL = res.SecureURL
	// saving to data base
	if err := u.repo.SaveDocument(ctx, doc); err != nil {
		u.log.Printf("ERR WHILE SAVING FILE URL TO DB::%v", err.Error())
		err = errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: "Failed to save file URL to database",
			Code:    http.StatusInternalServerError,
		}
		return err
	}

	return nil
}

func InitCloudinary() (*cloudinary.Cloudinary, error) {
	envFilePath := ".env"

	err := godotenv.Overload(envFilePath)
	if err != nil {
		log.Println("Error loading .env file:", err)
		return nil, errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	// Set the URL for the Cloudinary API
	fmt.Println("||||||||||||    CLOUDINARY VARIABLES    |||||||||||||||")
	log.Println("CLOUDINARY_CLOUD_NAME:", os.Getenv("CLOUDINARY_CLOUD_NAME"))
	log.Println("CLOUDINARY_API_KEY:", os.Getenv("CLOUDINARY_API_KEY"))
	log.Println("CLOUDINARY_API_SECRET:", os.Getenv("CLOUDINARY_API_SECRET"))

	if err != nil {
		log.Println("Error initializing Cloudinary:", err)
		return nil, errors.Error{
			Type:    "INTERNAL_ERROR",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	return cld, nil
}
