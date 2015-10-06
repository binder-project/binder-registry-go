package registry

import "encoding/json"

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

// Error implements the error interface and tries to return the serialized JSON
// form of the APIErrorResponse
func (apiError APIErrorResponse) Error() string {
	b, err := json.Marshal(apiError)
	if err != nil {
		return apiError.Message
	}
	return string(b)
}

var unavailableTemplateError = APIErrorResponse{Message: "Template unavailable"}
var existingTemplateError = APIErrorResponse{Message: "Template already exists"}
