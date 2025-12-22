package utils

import (
	"github.com/google/uuid"
)

// GenerateAPIKey generates a new API key using UUID
func GenerateAPIKey() string {
	return uuid.New().String()
}
