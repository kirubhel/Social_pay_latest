package usecase

import (
	"time"

	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"

	"github.com/google/uuid"
)

// GatewayInteractor handles gateway-related operations
type Interactor interface {
	// Gateway Management

	CreateGateway(
		name string,
		description string,
		gatewayType string,
		config entity.GatewayConfig,
		isActive bool,
	) (*entity.ListPaymentGateway, error)

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
	) (*entity.PaymentGateway, error)
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

// Main Interactor interface combining all capabilities
type GatewayInteractor interface {
}
