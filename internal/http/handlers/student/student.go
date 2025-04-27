package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/suryakantdsa/student-api/internal/storage"
	"github/suryakantdsa/student-api/internal/types"
	"github/suryakantdsa/student-api/internal/utils/response"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, err)
		}

		response.WriteResponse(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student details..!", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteResponse(w, http.StatusOK, student)

	}
}
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("$limit")
		skipStr := r.URL.Query().Get("$skip")
		name := r.URL.Query().Get("name")
		email := r.URL.Query().Get("email")
		age := r.URL.Query().Get("age")
		limit := 20
		skip := 0

		if limitStr != "" {
			if pasrsedLimit, err := strconv.Atoi(limitStr); err == nil {
				limit = pasrsedLimit
			}
		}

		if skipStr != "" {
			if pasrsedSkip, err := strconv.Atoi(skipStr); err == nil {
				skip = pasrsedSkip
			}
		}
		queryParams := map[string]string{
			"name":  name,
			"email": email,
			"age":   age,
		}
		slog.Info("getting all students")
		students, err := storage.GetStudents(limit, skip, queryParams)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteResponse(w, http.StatusOK, students)

	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		id := r.PathValue("id")
		intId, er := strconv.ParseInt(id, 10, 64)
		if er != nil {
			response.WriteResponse(w, http.StatusBadRequest, response.GeneralError(er))
			return
		}
		slog.Info("getting a student details..!", slog.String("id", id))

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

		updatedStudent, err := storage.UpdateStudent(intId, student)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, err)
		}

		response.WriteResponse(w, http.StatusOK, updatedStudent)

	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteResponse(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		message, err := storage.DeleteStudent(intId)
		if err != nil {
			response.WriteResponse(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteResponse(w, http.StatusOK, message)
	}
}
