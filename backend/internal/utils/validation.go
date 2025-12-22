package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleValidationError formats and returns validation errors
func HandleValidationError(c *gin.Context, err error) {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, formatValidationError(fieldError))
		}
	} else {
		errors = append(errors, err.Error())
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Validation failed",
		"details": errors,
	})
}

func formatValidationError(fieldError validator.FieldError) string {
	field := strings.ToLower(fieldError.Field())
	tag := fieldError.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "max":
		return field + " exceeds maximum length of " + fieldError.Param()
	case "min":
		return field + " must be at least " + fieldError.Param()
	default:
		return field + " is invalid"
	}
}
