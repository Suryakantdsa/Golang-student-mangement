package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/suryakantdsa/student-api/internal/types"
	"github/suryakantdsa/student-api/internal/utils/response"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			slog.Info("erro")
			response.WriteResponse(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validErr := err.(validator.ValidationErrors)
			response.WriteResponse(w, http.StatusBadRequest, response.ValidationError(validErr))
			return
		}

		response.WriteResponse(w, http.StatusCreated, student)
	}
}
