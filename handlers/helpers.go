package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/linn221/RequesterBackend/utils"
)

type Session struct{}

type MyHandlerFunc func(ctx context.Context, session Session, w http.ResponseWriter, r *http.Request) error

func Default(h MyHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mySession := Session{}
		ctx := r.Context()

		err := h(ctx, mySession, w, r)
		if err != nil {
			utils.RespondError(w, err)
		}
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// will parse the request, if found errors, will write to the response
// instance, continue, internalError
func writeValidationErrors(errs validator.ValidationErrors) error {
	var message string
	for _, err := range errs {
		message += fmt.Sprintf("'%s': %s\n", err.Field(), err.Tag())
	}
	return utils.BadRequest(message)
}

// will parse the request, if found errors, will write to the response
func parseJson[T any](r *http.Request) (*T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	err = validateStruct.Struct(&v)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			err := writeValidationErrors(ve)
			return nil, err
		}
		return nil, err
	}
	return &v, nil
}
