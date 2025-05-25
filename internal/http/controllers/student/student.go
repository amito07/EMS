package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/amito07/ems/internal/structure"
	"github.com/amito07/ems/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the EMS!"))
	}
}

func Create() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		var student structure.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF){
			response.WriteResponse(w, http.StatusBadRequest, "Empty body", response.GeneralErrorResponse(err))
			return
		}

		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, "Invalid request body", response.GeneralErrorResponse(err))
			return
		}

		// validate the request body
		if err := validator.New().Struct(student); err != nil {
			validteError := err.(validator.ValidationErrors)
			response.WriteResponse(w, http.StatusBadRequest, "Validation error", response.ValidationErrorResponse(validteError))
			return
		}


		response.WriteResponse(w, http.StatusCreated, "Student created successfully", student)
		
	}
}