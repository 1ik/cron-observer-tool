package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// validateUUID checks if the string is a valid UUID format
var validateUUID validator.Func = func(fl validator.FieldLevel) bool {
	uuidStr := fl.Field().String()
	if uuidStr == "" {
		return true // Let required tag handle empty values
	}
	_, err := uuid.Parse(uuidStr)
	return err == nil
}
