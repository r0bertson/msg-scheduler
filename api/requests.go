package api

import (
	"errors"
	"github.com/badoux/checkmail"
	"msg-scheduler/common/messaging"
	"msg-scheduler/common/models"
	"strings"
)

type UserOperationsRequestBody struct {
	Email    string `binding:"email"`
	Password string
}

type MessageOperationsRequestBody struct {
	Subject string
	Content string
}

type SendEmailRequestBody struct {
	Timeout *int
	Payload messaging.EmailPayload
}

func (op *UserOperationsRequestBody) Validate(action models.Action) error {
	switch action {
	case models.Create:
		if op.Password == "" {
			return fieldRequiredError("Password")
		}
		if op.Email == "" {
			return fieldRequiredError("Email")
		}
		if err := checkmail.ValidateFormat(op.Email); err != nil {
			return errors.New("invalid email")
		}

		return nil
	case models.Update:
		if op.Password == "" {
			return fieldCannotBeEmptyError("Password")
		}
		if op.Email == "" {
			return fieldCannotBeEmptyError("Email")
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
			return fieldRequiredError("Content")
		}
		if strings.TrimSpace(op.Subject) == "" {
			return fieldRequiredError("Subject")
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
