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

// Error implements the error interface and tries to return the serialized JSON
// form of the APIErrorResponse
func (apiError APIErrorResponse) Error() string {
	return apiError.Message
}

// UnavailableTemplateError is an APIErrorResponse with a boilerplate unavailable template message
var UnavailableTemplateError = APIErrorResponse{Message: "Template unavailable"}

// ExistingTemplateError is an APIErrorResponse with a boilerplate message for when a template already exists
var ExistingTemplateError = APIErrorResponse{Message: "Template already exists"}

// DontPanicError is an APIErrorResponse when there's a mostly catastrophic error
var DontPanicError = APIErrorResponse{Message: "Internal Server Error. Don't Panic. We will."}

// UnableToListError is an APIErrorResponse reported when templates are unable to be listed
var UnableToListError = APIErrorResponse{Message: "Unable to list templates"}

// TemplateNotFoundError is an APIErrorResponse reported when the template requested doesn't exist
var TemplateNotFoundError = APIErrorResponse{Message: "Template Not Found"}

var InvalidParameterError = APIErrorResponse{Message: "Invalid parameters provided"}
