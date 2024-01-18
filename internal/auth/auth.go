package auth

type ErrorResponse struct {
	Code    uint16   `json:"error_code"`
	Message string   `json:"error_message"`
	Details []string `json:"error_details"`
}

// Error returns the error message for the ErrorResponse
func (err ErrorResponse) Error() string {
	return err.Message
}

type User struct {
	Id           string   `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Verified     bool     `json:"verified"`
	Type         uint8    `json:"type"`
	Provider     string   `json:"provider"`
	CreatedAtUTC string   `json:"created_at_utc"`
	Features     []string `json:"features"`
}

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	UserId       string   `json:"user_id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Verified     bool     `json:"verified"`
	Provider     string   `json:"provider"`
	CreatedAtUTC string   `json:"created_at_utc"`
	Features     []string `json:"features"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token         string `json:"token"`
	TokenType     string `json:"token_type"`
	ApplicationId string `json:"application"`
	ExpiresIn     int64  `json:"expires_in"`
	UserId        string `json:"user_id"`
	Username      string `json:"username"`
	RefreshToken  string `json:"refresh_token"`
}

type UserAccountVerificationRequest struct {
	UserId            string `json:"user_id"`
	VerificationToken string `json:"verification_token"`
}

type UserPasswordResetInitiationRequest struct {
	Email string `json:"email"`
}

type UserPasswordResetExecutionRequest struct {
	// The user id of the user to reset the password for
	UserID string `json:"user_id"`
	// The new password for the user
	NewPassword string `json:"new_password"`
	// The reset token for the user
	PasswordResetToken string `json:"password_reset_token"`
	// The verification code for the user
	VerificationCode string `json:"verification_code"`
}

const (
	// Error codes
	// Error code 0 indicates an unhandled error. This means there was a server error.
	UnhandledError        = 1
	UnhandledErrorMessage = "an unhandled/unexpected error occured"
	// Error code 5 indicates the request body could not be parsed or was otherwise invalid.
	RequestPayloadInvalid     = 5
	RequestBodyInvalidMessage = "the request body could not be parsed"
	// Error code 10 indicates the request failed validation.
	// This means the request content was parsed but failed validation of the content.
	RequestValidationFailure        = 10
	RequestValidationFailureMessage = "request validation failure"
	// Error code 15 indicates the requested application resource was not found.
	ApplicationNotFound        = 15
	ApplicationNotFoundMessage = "application not found"
	// Error code 20 indicates the credentials provided were invalid.
	InvalidCredentials        = 20
	InvalidCredentialsMessage = "invalid credentials"
	// Error code 25 indicates the data provided conflicts with existing data.
	// This means that the data provided cannot be used because it conflicts with existing data.
	DataConflict        = 25
	DataConflictMessage = "data conflict"
	// Error code 30 indicates the user has not verified their email address.
	UserNotVerified        = 30
	UserNotVerifiedMessage = "user not verified"
	// Error code 35 indicates the provided authorization token was invalid or has been blacklisted.
	InvalidAuthToken        = 35
	InvalidAuthTokenMessage = "invalid or malformed authorization token"
	// Error code 40 indicates the user does not have access to the requested resource.
	AccessDenied        = 40
	AccessDeniedMessage = "access denied"
	// Error code 45 indicates the provided verification code was invalid.
	InvalidUserVerficationToken        = 45
	InvalidUserVerficationTokenMessage = "invalid verification code"
	// Error code 50 indicates the requested user resource was not found.
	UserNotFound        = 50
	UserNotFoundMessage = "user not found"
	// Error code 55 indicates the provided password reset token was invalid.
	InvalidPasswordResetToken        = 55
	InvalidPasswordResetTokenMessage = "invalid password reset token"
	// Error code 60 indicates the provided password reset verification code was invalid.
	InvalidPasswordResetVerificationCode        = 60
	InvalidPasswordResetVerificationCodeMessage = "invalid password reset verification code"
	// Error code 65 indicates that missing or invalid request headers were provided.
	InvalidRequestHeaders        = 65
	InvalidRequestHeadersMessage = "invalid or missing request headers"
	// Error code 70 indicates that an authorization token is expired.
	AuthTokenExpired        = 70
	AuthTokenExpiredMessage = "authorization token expired"
)
