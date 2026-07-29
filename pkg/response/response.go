package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ErrMessage struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

var (
	InternalServerErrResponse = ErrMessage{
		Message: "Processing error",
		Details: nil,
	}
	InputErrResponse = ErrMessage{
		Message: "Invalid input",
		Details: nil,
	}
)

// InputFieldErrorResponse parses a validation error and formats it into an ErrMessage response.
func InputFieldErrorResponse(err error) ErrMessage {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}

	var errs []string
	for _, err:= range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field() + " is invalid (" + err.Tag() + ")")
	}
	return ErrMessage{
		Message: "Invalid input",
		Details: errs,
	}
}