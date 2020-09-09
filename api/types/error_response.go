package types

type ErrorResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Title    string `json:"message_title"`
	Severity string `json:"message_severity"`
}
