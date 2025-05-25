package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// const valid status codes
const (
	StatusOk = "OK"
	StatusError = "ERROR"
)

type Response struct {
	Status string `json:"status"` // Status of the response
	Error string `json:"error"` // Error message if any
}

func WriteResponse(w http.ResponseWriter, statusCode int, message string, data any) error{
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]any{
		"status":  statusCode,
		"message": message,
		"data":    data,
	}

	return json.NewEncoder(w).Encode(response);
}

//Generic function for error response
func GeneralErrorResponse(err error) Response {
	fmt.Println("Error........")
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationErrorResponse(errs validator.ValidationErrors) Response {
	var errorMessages []string
	for _, err := range errs {
		switch err.ActualTag(){
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is required", err.Field()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation for tag '%s'", err.Field(), err.ActualTag()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errorMessages, ", "),
	}
}