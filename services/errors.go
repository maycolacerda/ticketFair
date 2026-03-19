// services/errors.go
package services

import "errors"

var (
	// Auth
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrMerchantDisabled   = errors.New("merchant account is disabled")
	ErrUnauthorized       = errors.New("unauthorized")

	// Not found
	ErrUserNotFound        = errors.New("user not found")
	ErrMerchantNotFound    = errors.New("merchant not found")
	ErrRepNotFound         = errors.New("merchant representative not found")
	ErrEventNotFound       = errors.New("event not found")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrProfileNotFound     = errors.New("profile not found")

	// Conflict
	ErrEmailInUse    = errors.New("email already in use")
	ErrUsernameInUse = errors.New("username already in use")
	ErrPhoneInUse    = errors.New("phone number already in use")
	ErrProfileExists = errors.New("profile already exists")

	// Validation
	ErrNoFieldsToUpdate = errors.New("no fields to update")
	ErrInvalidTimeRange = errors.New("end_time must be after start_time")
	ErrStartTimeInPast  = errors.New("start_time must be in the future")
	ErrInvalidCapacity  = errors.New("capacity must be greater than zero")

	// Tickets
	ErrEventSoldOut  = errors.New("event is sold out")
	ErrNotRefundable = errors.New("transaction is not refundable")

	// Internal
	ErrFailedToCreate        = errors.New("failed to create record")
	ErrFailedToUpdate        = errors.New("failed to update record")
	ErrFailedToFetch         = errors.New("failed to fetch record")
	ErrFailedToHash          = errors.New("failed to process password")
	ErrFailedToGenerateToken = errors.New("failed to generate token")
)
