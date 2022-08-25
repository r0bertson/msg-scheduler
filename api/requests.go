package api

import (
	"errors"
	"github.com/badoux/checkmail"
	"msg-scheduler/common/models"
	"strings"
)

type UserOperationsRequestBody struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password"`
}

type MessageOperationsRequestBody struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type TestMessageOperationsRequestBody struct {
	To      string `json:"to" binding:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (op *UserOperationsRequestBody) Validate(action models.Action) error {
	switch action {
	case models.Create:
		if op.Password == "" {
			return fieldRequiredError("password")
		}
		if op.Email == "" {
			return fieldRequiredError("email")
		}
		if err := checkmail.ValidateFormat(op.Email); err != nil {
			return errors.New("invalid email")
		}

		return nil
	case models.Update:
		if op.Password == "" {
			return fieldCannotBeEmptyError("password")
		}
		if op.Email == "" {
			return fieldCannotBeEmptyError("email")
		}
		if err := checkmail.ValidateFormat(op.Email); err != nil {
			return errors.New("invalid email format")
		}
		return nil
	default:
		return nil
	}
}

func (op *MessageOperationsRequestBody) Validate(action models.Action) error {
	switch action {
	case models.Create, models.Update:
		if strings.TrimSpace(op.Content) == "" {
			return fieldRequiredError("content")
		}
		if strings.TrimSpace(op.Subject) == "" {
			return fieldRequiredError("subject")
		}
	}

	return nil
}

func fieldRequiredError(field string) error {
	return errors.New("required field: " + field)
}

func fieldCannotBeEmptyError(field string) error {
	return errors.New(field + " cannot be empty")
}
