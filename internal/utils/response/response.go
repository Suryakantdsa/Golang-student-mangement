package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `jsom:"error"`
}

const (
	statusOk  = "Ok"
	statusErr = "Error"
)

func WriteResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: statusErr,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s filed is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s filed is invalid", err.Field()))

		}

	}
	return Response{
		Status: statusErr,
		Error:  strings.Join(errMsgs, ","),
	}

}
