package config

import (
	"fmt"
	"strings"
)

// Validate checks that all required configuration fields are set
func (c *Config) Validate() error {
	var missing []string

	// Check required database fields
	if c.Database.URI == "" {
		missing = append(missing, "DATABASE_URI")
	}
	if c.Database.Name == "" {
		missing = append(missing, "DATABASE_NAME")
	}

	if len(missing) > 0 {
		return &MissingConfigError{Fields: missing}
	}

	return nil
}

// MissingConfigError represents missing required configuration
type MissingConfigError struct {
	Fields []string
}

func (e *MissingConfigError) Error() string {
	var sb strings.Builder

	sb.WriteString("❌ Missing Required Configuration\n\n")
	sb.WriteString("The following environment variables are required but not set:\n")

	for _, field := range e.Fields {
		sb.WriteString(fmt.Sprintf("  • %s\n", field))
	}

	sb.WriteString("\nTo fix this, either:\n")
	sb.WriteString("1. Set environment variables:\n")
	for _, field := range e.Fields {
		sb.WriteString(fmt.Sprintf("   export %s=your_value\n", field))
	}

	sb.WriteString("\n2. Create a .env file in the backend directory:\n")
	for _, field := range e.Fields {
		sb.WriteString(fmt.Sprintf("   %s=your_value\n", field))
	}

	sb.WriteString("\n3. Or pass them when running:\n")
	sb.WriteString("   ")
	for i, field := range e.Fields {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%s=your_value", field))
	}
	sb.WriteString(" go run cmd/server/main.go")

	return sb.String()
}


