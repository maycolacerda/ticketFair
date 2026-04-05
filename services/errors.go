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

	//verification
	ErrVerificationNotFound = errors.New("verification code not found")
	ErrVerificationExpired  = errors.New("verification code expired")
	ErrVerificationUsed     = errors.New("verification code already used")
	ErrVerificationInvalid  = errors.New("invalid verification code")
	ErrAlreadyVerified      = errors.New("already verified")
	ErrProfilePhoneNotFound = errors.New("phone number not found on profile")

	//admin
	ErrAdminNotFound = errors.New("admin not found")
	ErrAdminDisabled = errors.New("admin account is disabled")

	//s3
	ErrS3NotConfigured  = errors.New("image storage is not configured")
	ErrImageTooLarge    = errors.New("image exceeds maximum size of 5MB")
	ErrInvalidImageType = errors.New("only JPEG, PNG and WebP images are allowed")

	//password recovery
	ErrResetCodeNotFound = errors.New("reset code not found")
	ErrResetCodeExpired  = errors.New("reset code has expired")
	ErrResetCodeUsed     = errors.New("reset code already used")
	ErrResetCodeInvalid  = errors.New("invalid reset code")

	//ticket type
	ErrTicketTypeNotFound   = errors.New("ticket type not found")
	ErrTicketTypeSoldOut    = errors.New("ticket type is sold out")
	ErrTicketSaleNotStarted = errors.New("ticket sales have not started yet")
	ErrTicketSaleEnded      = errors.New("ticket sales have ended")
	ErrTicketBelowMinimum   = errors.New("quantity is below the minimum per order")
	ErrTicketExceedsMaximum = errors.New("quantity exceeds the maximum per order")
	ErrTicketTypeHasSales   = errors.New("cannot delete ticket type with existing sales")

	ErrConnectionNotFound   = errors.New("connection not found")
	ErrConnectionExists     = errors.New("connection already exists")
	ErrConnectionSelf       = errors.New("cannot connect with yourself")
	ErrNotConnected         = errors.New("users are not connected")
	ErrConnectionNotPending = errors.New("connection is not pending")
	ErrTicketAlreadyGifted  = errors.New("ticket has already been gifted")
	ErrCannotGiftOwn        = errors.New("cannot gift a ticket to yourself")
	ErrTicketNotActive      = errors.New("ticket is not active")
)
