package entity

// Resource represents different resources in the system
type Resource string

const (
	RESOURCE_ALL          Resource = "ALL"       // Special permission for merchants to access all resources
	RESOURCE_ADMIN_ALL    Resource = "ADMIN_ALL" // Special permission for super admins to access all admin resources
	RESOURCE_TRANSACTION  Resource = "transaction"
	RESOURCE_MERCHANT     Resource = "merchant"
	RESOURCE_USER         Resource = "user"
	RESOURCE_ADMIN_WALLET Resource = "admin_wallet"
	RESOURCE_IP_WHITELIST Resource = "ip_whitelist"
	RESOURCE_API_KEY      Resource = "api_key"
	RESOURCE_WEBHOOK      Resource = "webhook"
	RESOURCE_ANALYTICS    Resource = "analytics"
	RESOURCE_COMMISSION   Resource = "commission"
	RESOURCE_QR           Resource = "qr"
	RESOURCE_CHECKOUT     Resource = "checkout"
	RESOURCE_NOTIFICATION Resource = "notification"
	RESOURCE_WALLET       Resource = "wallet"
)

// Operation represents different operations that can be performed
type Operation string

const (
	OPERATION_ALL          Operation = "ALL"       // Special permission for merchants to perform all operations
	OPERATION_ADMIN_ALL    Operation = "ADMIN_ALL" // Special permission for super admins to perform all admin operations
	OPERATION_CREATE       Operation = "CREATE"
	OPERATION_READ         Operation = "READ"
	OPERATION_UPDATE       Operation = "UPDATE"
	OPERATION_DELETE       Operation = "DELETE"
	OPERATION_ADMIN_CREATE Operation = "ADMIN_CREATE"
	OPERATION_ADMIN_READ   Operation = "ADMIN_READ"
	OPERATION_ADMIN_UPDATE Operation = "ADMIN_UPDATE"
	OPERATION_ADMIN_DELETE Operation = "ADMIN_DELETE"
)

// UserType represents user types in the system
type UserType string

const (
	USER_TYPE_SUPER_ADMIN UserType = "super_admin"
	USER_TYPE_ADMIN       UserType = "admin"
	USER_TYPE_MERCHANT    UserType = "merchant"
	USER_TYPE_MEMBER      UserType = "member"
)

// Gender represents user gender
type Gender string

const (
	GENDER_MALE   Gender = "MALE"
	GENDER_FEMALE Gender = "FEMALE"
)
