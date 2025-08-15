package seeder

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/repository"
)

// MerchantSeeder handles seeding merchant
type MerchantSeeder struct {
	db     *sql.DB
	logger *log.Logger
	repo   repository.AuthRepository
}

// FakeMerchantDatas provide type definition to generate data using faker
type FakeMerchantDatas struct {
	// user data
	Title     string `faker:"title_male"`
	FirstName string `faker:"first_name"`
	LastName  string `faker:"last_name"`
	Email     string `faker:"email"`
	Password  string `faker:"password"`

	// merchant data
	LegalName                  string `faker:"name"`
	TradingName                string `faker:"name"`
	BusinessRegistrationNumber int64  `faker:"unix_time"`
	TaxID                      int64  `faker:"unix_time"`
	LotteryCertificateNumber   int64  `faker:"unix_time"`
	WebsiteURL                 string `faker:"url"`
	EstablishedDate            int64  `faker:"unix_time"`
}

// NewMerchantSeeder creates new merchant seeder
func NewMerchantSeeder(db *sql.DB, repo repository.AuthRepository) *MerchantSeeder {
	return &MerchantSeeder{
		db:     db,
		logger: log.New(os.Stdout, "", log.LstdFlags),
		repo:   repo,
	}
}

// SeedMerchant seeds n number of merchants
func (s *MerchantSeeder) SeedMerchant(ctx context.Context, n int) error {
	for i := range n {
		// generate fake data using faker
		fakeMerchant := FakeMerchantDatas{}
		err := faker.FakeData(&fakeMerchant)

		if err != nil {
			s.logger.Printf("Error generating fake data: %v", err)
		}

		// seeding new user data
		s.logger.Printf("Seeding the %v merchant", i+1)
		randomPhone := rand.Intn(888888888) + 111111111
		user, err := s.repo.CreateUser(ctx, &entity.CreateUserRequest{
			Title:        fakeMerchant.Title,
			FirstName:    fakeMerchant.FirstName,
			LastName:     fakeMerchant.LastName,
			Email:        fakeMerchant.Email,
			PhonePrefix:  "251",
			PhoneNumber:  strconv.Itoa(randomPhone),
			Password:     fakeMerchant.Password,
			PasswordHint: "",
			UserType:     "merchant",
		})

		if err != nil {
			s.logger.Printf("Error seeding user: %v", err)
		}

		// seeding new merchant data
		merchant, err := s.repo.CreateMerchant(ctx, user.ID, map[string]interface{}{
			"legal_name":                   fakeMerchant.LegalName,
			"trading_name":                 fakeMerchant.TradingName,
			"business_type":                "Individual",
			"business_registration_number": strconv.FormatInt(fakeMerchant.BusinessRegistrationNumber, 10),
			"tax_identification_number":    strconv.FormatInt(fakeMerchant.TaxID, 10),
			"industry_category":            "Finance",
			"is_betting_company":           true,
			"lottery_certificate_number":   strconv.FormatInt(fakeMerchant.LotteryCertificateNumber, 10),
			"website_url":                  fakeMerchant.WebsiteURL,
			"established_date":             time.Unix(fakeMerchant.EstablishedDate, 0),
			"status":                       "active",
		})

		if err != nil {
			s.logger.Printf("Error seeding merchant: %v", err)
		}

		// seeding new merchant contact data
		id := uuid.New()
		_, err = s.db.ExecContext(ctx, `
			INSERT INTO merchants.contacts (id, merchant_id, contact_type, first_name, last_name, email, phone_number, "position", is_verified, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
			`, id, merchant, "", fakeMerchant.FirstName, fakeMerchant.LastName, fakeMerchant.Email, user.PhoneNumber, "", true, time.Now(), time.Now())

		if err != nil {
			s.logger.Printf("Error seeding merchant contact: %v", err)
		}

		// seeding new merchant document data
		id = uuid.New()
		_, err = s.db.ExecContext(ctx, `
			INSERT INTO merchants.documents
			(id, merchant_id, document_type, document_number, file_url, file_hash, verified_by, verified_at, status, rejection_reason, created_at, updated_at)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
			`, id, merchant, "ID", "123", "", "", nil, nil, "pending", nil, time.Now(), time.Now())

		if err != nil {
			s.logger.Printf("Error seeding merchant document: %v", err)
		}
	}

	return nil
}
