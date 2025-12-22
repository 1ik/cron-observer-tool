package validators

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// validateObjectID checks if the string is a valid MongoDB ObjectID format
var validateObjectID validator.Func = func(fl validator.FieldLevel) bool {
	objectIDStr := fl.Field().String()
	if objectIDStr == "" {
		return true // Let required tag handle empty values
	}
	_, err := primitive.ObjectIDFromHex(objectIDStr)
	return err == nil
}

// validateCron checks if the string is a valid cron expression
var validateCron validator.Func = func(fl validator.FieldLevel) bool {
	cronStr := fl.Field().String()
	if cronStr == "" {
		return true // Let required tag handle empty values
	}

	// Basic cron expression validation: 5 fields (minute hour day month weekday)
	// or 6 fields (second minute hour day month weekday)
	parts := strings.Fields(cronStr)
	if len(parts) != 5 && len(parts) != 6 {
		return false
	}

	// Validate each part contains valid cron characters
	cronPattern := regexp.MustCompile(`^[\d\*\-\,\/]+$`)
	for _, part := range parts {
		if !cronPattern.MatchString(part) {
			return false
		}
	}

	return true
}

// validateTimezone checks if the string is a valid timezone
var validateTimezone validator.Func = func(fl validator.FieldLevel) bool {
	timezoneStr := fl.Field().String()
	if timezoneStr == "" {
		return true // Let required tag handle empty values
	}
	_, err := time.LoadLocation(timezoneStr)
	return err == nil
}

// validateTimeFormat checks if the string is in HH:MM format
var validateTimeFormat validator.Func = func(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()
	if timeStr == "" {
		return true // Let required tag handle empty values
	}

	// Validate HH:MM format
	timePattern := regexp.MustCompile(`^([0-1][0-9]|2[0-3]):[0-5][0-9]$`)
	if !timePattern.MatchString(timeStr) {
		return false
	}

	// Try parsing to ensure it's valid
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// validateURL checks if the string is a valid URL format
var validateURL validator.Func = func(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()
	if urlStr == "" {
		return true // Let required tag handle empty values
	}
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

// validateHTTPMethod checks if the string is a valid HTTP method
var validateHTTPMethod validator.Func = func(fl validator.FieldLevel) bool {
	method := strings.ToUpper(fl.Field().String())
	if method == "" {
		return true // Let required tag handle empty values
	}
	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
	}
	return validMethods[method]
}

// RegisterCustomValidators registers all custom validators with the validator instance
func RegisterCustomValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("uuid", validateUUID); err != nil {
		return err
	}
	if err := v.RegisterValidation("objectid", validateObjectID); err != nil {
		return err
	}
	if err := v.RegisterValidation("cron", validateCron); err != nil {
		return err
	}
	if err := v.RegisterValidation("timezone", validateTimezone); err != nil {
		return err
	}
	if err := v.RegisterValidation("time_format", validateTimeFormat); err != nil {
		return err
	}
	if err := v.RegisterValidation("url", validateURL); err != nil {
		return err
	}
	if err := v.RegisterValidation("http_method", validateHTTPMethod); err != nil {
		return err
	}
	return nil
}
