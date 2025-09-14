package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

var (
	ErrNotFound   = errors.New("record not found")
	ErrBadRequest = errors.New("bad request")
	ErrDb         = errors.New("db error")
)

func RespondError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, ErrNotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, ErrBadRequest) {
		status = http.StatusBadRequest
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	}
	http.Error(w, err.Error(), status)
}

func OkCreated(w http.ResponseWriter, id int) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, id)
}

func OkUpdated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func OkDeleted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(message string) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, message)
}

func GetIdParam(r *http.Request) (int, error) {

	resIdStr := r.PathValue("id")
	resId, err := strconv.Atoi(resIdStr)
	if err != nil || resId <= 0 {
		return 0, BadRequest("invalid resource id")
	}

	return resId, nil
}

func OkJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	FinalErrHandle(w, json.NewEncoder(w).Encode(v))
}

func FinalErrHandle(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
