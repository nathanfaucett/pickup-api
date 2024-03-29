package model

import "github.com/gofiber/fiber/v2"

type ErrorMessageST struct {
	Message    string        `json:"error" validate:"required"`
	Parameters []interface{} `json:"parameters" validate:"required"`
} // @name ErrorMessage

func NewErrorMessage(message string, parameters ...interface{}) *ErrorMessageST {
	return &ErrorMessageST{
		Message:    message,
		Parameters: parameters,
	}
}

type ErrorST struct {
	StatusCode int                          `json:"-"`
	Errors     map[string][]*ErrorMessageST `json:"errors" validate:"required"`
} // @name Error

func NewError(statusCode int) *ErrorST {
	return &ErrorST{
		StatusCode: statusCode,
		Errors:     make(map[string][]*ErrorMessageST),
	}
}

func (e *ErrorST) AddError(name string, message string, parameters ...interface{}) *ErrorST {
	if e.Errors[name] == nil {
		e.Errors[name] = make([]*ErrorMessageST, 0)
	}
	e.Errors[name] = append(e.Errors[name], NewErrorMessage(message, parameters...))
	return e
}

func (e *ErrorST) Send(c *fiber.Ctx) error {
	return c.Status(e.StatusCode).JSON(e)
}
