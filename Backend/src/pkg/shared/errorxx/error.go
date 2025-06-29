package errorxx

import "github.com/joomcode/errorx"

// Top-level namespaces
var (
	AppError      = errorx.NewNamespace("APPLICATION::ERR")      // User-caused, client-side errors
	AuthError     = errorx.NewNamespace("AUTHENTICATION::ERR")   // Authentication/authorization failures
	BusinessError = errorx.NewNamespace("BUSINESS_LOGIC::ERR")   // Domain/business logic violations
	FileError     = errorx.NewNamespace("FILE_IO::ERR")          // File/media upload/download issues
	ExternalError = errorx.NewNamespace("EXTERNAL_SERVICE::ERR") // API, 3rd-party service failures
	InternalError = errorx.NewNamespace("INTERNAL_SERVER::ERR")  // Unexpected server-side/internal issues
	DatabaseError = errorx.NewNamespace("DATABASE::ERR")         // Persistence layer problems
)

// Application (client-side) errors
var (
	ErrAppBadInput      = errorx.NewType(AppError, "validation_bad_input")      // Invalid user input
	ErrAppMissingField  = errorx.NewType(AppError, "validation_missing_field")  // Required field is missing
	ErrAppInvalidFormat = errorx.NewType(AppError, "validation_invalid_format") // Wrong data format (e.g., email)
)

// Authentication & Authorization errors
var (
	ErrAuthUnauthorized = errorx.NewType(AuthError, "unauthorized")  // Not logged in
	ErrAuthForbidden    = errorx.NewType(AuthError, "forbidden")     // Logged in, but not allowed
	ErrAuthTokenExpired = errorx.NewType(AuthError, "token_expired") // Expired JWT or session
)

// Business logic errors
var (
	ErrBizConflict      = errorx.NewType(BusinessError, "business_conflict")  // Conflict in domain rules
	ErrBizNotAllowed    = errorx.NewType(BusinessError, "action_not_allowed") // Action not allowed in current state
	ErrBizLimitExceeded = errorx.NewType(BusinessError, "limit_exceeded")     // User quota/limit exceeded
)

// File or media handling errors
var (
	ErrFileTooLarge    = errorx.NewType(FileError, "file_too_large")        // Upload exceeds size limit
	ErrFileUnsupported = errorx.NewType(FileError, "unsupported_file_type") // File type not allowed
	ErrFileReadFail    = errorx.NewType(FileError, "file_read_error")       // Failed to read file
)

// External services (3rd-party APIs, webhooks)
var (
	ErrExternalTimeout     = errorx.NewType(ExternalError, "service_timeout")     // Timed out waiting for external API
	ErrExternalInvalidRes  = errorx.NewType(ExternalError, "invalid_response")    // Malformed or unexpected response
	ErrExternalUnavailable = errorx.NewType(ExternalError, "service_unavailable") // Remote service is down
)

// Internal system/runtime errors
var (
	ErrInternalMarshal    = errorx.NewType(InternalError, "marshal_failed")   // Failed to encode data (e.g., JSON)
	ErrInternalUnmarshal  = errorx.NewType(InternalError, "unmarshal_failed") // Failed to decode data
	ErrInternalPanic      = errorx.NewType(InternalError, "panic_recovered")  // Panic caught and handled
	ErrInternalUnexpected = errorx.NewType(InternalError, "unexpected_state") // Should never happen; logic bug
)

// Database errors
var (
	ErrDBWrite      = errorx.NewType(DatabaseError, "db_write_failed")   // Insert/update failed
	ErrDBRead       = errorx.NewType(DatabaseError, "db_read_failed")    // Query/select failed
	ErrDBNullID     = errorx.NewType(DatabaseError, "null_object_id")    // Missing expected object ID
	ErrDBDuplicate  = errorx.NewType(DatabaseError, "duplicate_record")  // Unique constraint violation
	ErrDBConnection = errorx.NewType(DatabaseError, "connection_failed") // DB unreachable or timed out
)

// Custom error property (optional tagging)
var (
	ErrorCode = errorx.RegisterProperty("ERRCODE")
)
