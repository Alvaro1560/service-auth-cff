package validation_email

type VerificationRequest struct {
	Email string `json:"email"`
}
type VerificationDataRequest struct {
	Id             int64  `json:"id"`
	Identification string `json:"identification"`
	Code           string `json:"code"`
}
