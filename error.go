package main

import (
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ValidationErrorResponse represents the structure of validation errors
type ValidationErrorResponse struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors"`
}

// FieldError represents a single field validation error
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

// GlobalErrorHandler is a centralized error handler for Echo
func GlobalErrorHandler(err error, c echo.Context) {
	// Don't handle if response already started
	if c.Response().Committed {
		return
	}

	// Handle validation errors
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		HandleValidationError(c, validationErrs)
		return
	}

	// Handle Echo HTTP errors
	if he, ok := err.(*echo.HTTPError); ok {
		code := he.Code
		message := he.Message

		// If message is a string, use it directly
		if msg, ok := message.(string); ok {
			c.JSON(code, map[string]interface{}{
				"message": msg,
			})
			return
		}

		// Otherwise use the message as is
		c.JSON(code, message)
		return
	}

	// Default to 500 Internal Server Error
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"message": "Internal server error",
		"error":   err.Error(),
	})
}

// HandleValidationError converts validator errors to 422 response
func HandleValidationError(c echo.Context, validationErrs validator.ValidationErrors) {
	errors := make(map[string]interface{})

	for _, err := range validationErrs {
		fieldName := err.Field()
		errors[fieldName] = FieldError{
			Field:   fieldName,
			Message: getValidationMessage(err),
			Tag:     err.Tag(),
			Value:   getValueString(err.Value()),
		}
	}

	response := ValidationErrorResponse{
		Message: "Validation failed",
		Errors:  errors,
	}

	c.JSON(http.StatusUnprocessableEntity, response)
}

// getValidationMessage returns a human-readable error message based on the validation tag
func getValidationMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "url":
		return field + " must be a valid URL"
	case "uuid":
		return field + " must be a valid UUID"
	case "min":
		return field + " must be at least " + err.Param()
	case "max":
		return field + " must be at most " + err.Param()
	case "gte":
		return field + " must be greater than or equal to " + err.Param()
	case "lte":
		return field + " must be less than or equal to " + err.Param()
	case "gt":
		return field + " must be greater than " + err.Param()
	case "lt":
		return field + " must be less than " + err.Param()
	case "len":
		return field + " must be " + err.Param() + " characters long"
	case "oneof":
		return field + " must be one of: " + err.Param()
	case "required_if":
		return field + " is required when " + err.Param()
	default:
		return field + " failed validation on " + err.Tag()
	}
}

// getValueString converts the field value to a string representation
func getValueString(value interface{}) string {
	if value == nil {
		return ""
	}

	// Don't expose sensitive values
	return ""
}

type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements the echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you can customize the error response here
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator returns a singleton validator instance
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
	return validate
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}
