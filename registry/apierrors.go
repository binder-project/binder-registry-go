package registry

// ErrorHint allows you to specify specific fields the user is potentially
// missing
type ErrorHint struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

// APIErrorResponse mirrors the GitHub error model, providing a message as well
// as possible error hints
type APIErrorResponse struct {
	Message string      `json:"message"`
	Errors  []ErrorHint `json:"errors,omitempty"`
}
