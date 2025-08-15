package usecase

import (
	"time"

	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"

	"github.com/google/uuid"
)

// AuthNRepo interface
type AuthNRepo interface {
	// Gateway Management

	CreateGateway(gateway entity.ListPaymentGateway) error

	GetGatewayByID(gatewayID uuid.UUID,
		includeInactive bool,
		limit int,
		offset int) ([]entity.GatewayMerchant, int, error)
	ListAllGateways() ([]entity.PaymentGateway, error)
	UpdateGateway(
		id uuid.UUID,
		name *string,
		description *string,
		isActive *bool,
		config *entity.GatewayConfig,
	) (*entity.PaymentGateway, error)
	DeleteGateway(id uuid.UUID) error

	// Merchant-Gateway Association
	LinkGatewayToMerchant(
		merchantID uuid.UUID,
		gatewayID uuid.UUID,
	) (time.Time, error)

	UnlinkGatewayFromMerchant(
		merchantID uuid.UUID,
		gatewayID uuid.UUID,
	) (time.Time, error)

	DisableMerchantGateway(
		merchantID uuid.UUID,
		gatewayID uuid.UUID,
		reason string,
	) (time.Time, error)

	GetGatewayMerchants(
		gatewayID uuid.UUID,
		includeInactive bool,
		limit, offset int,
	) (*entity.GatewayMerchantsResult, error)
	ListMerchantGateways(
		merchantID uuid.UUID,
		includeDisabled bool,
		gatewayType string,
	) ([]entity.GatewayMerchant, error)
	EnableMerchantGateway(
		merchantID uuid.UUID,
		gatewayID uuid.UUID,
		reason string) (time.Time, error)
}

// Repo interface, extending AuthNRepo
type Repo interface {
	AuthNRepo
}
