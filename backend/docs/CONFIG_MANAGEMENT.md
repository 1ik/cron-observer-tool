# Configuration Management with Viper

## Overview

This backend uses [Viper](https://github.com/spf13/viper) for centralized configuration management. Viper provides a flexible, type-safe way to manage application configuration from multiple sources with automatic precedence handling.

## Architecture

```mermaid
flowchart TD
    ENV[".env file"] --> Viper
    EnvVars["OS Environment Variables"] --> Viper
    Defaults["Code Defaults"] --> Viper
    Viper --> ConfigStruct["Config Struct"]
    ConfigStruct --> Validate["Validation"]
    Validate -->|Pass| Main["main.go"]
    Validate -->|Fail| Exit["Exit with Error"]
    Main --> Database["database.NewConnection(cfg.Database)"]
    Main --> Server["http.Server with cfg.Server"]
```

## Configuration Priority

Configuration values are loaded in the following priority order (highest to lowest):

1. **OS Environment Variables** (highest priority)
2. **`.env` file** (in backend directory)
3. **Code Defaults** (lowest priority)

This means environment variables will always override `.env` file values, which will override defaults.

## Quick Start

### 1. Create `.env` File

Create a `.env` file in the `backend/` directory:

```bash
# Required
DATABASE_URI=mongodb://localhost:27017
DATABASE_NAME=cronobserver

# Optional (defaults shown)
SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
DATABASE_TIMEOUT=10s
DATABASE_MAX_CONNS=100
```

### 2. Run the Application

```bash
cd backend
go run cmd/server/main.go
```

The application will:

- Load configuration from `.env` file
- Validate required fields
- Fail fast with clear error messages if required config is missing

## Configuration Values

### Required Fields

These fields **must** be provided (via `.env` file or environment variables):

| Config Key      | Environment Variable | Description            | Example                     |
| --------------- | -------------------- | ---------------------- | --------------------------- |
| `database.uri`  | `DATABASE_URI`       | MongoDB connection URI | `mongodb://localhost:27017` |
| `database.name` | `DATABASE_NAME`      | Database name          | `cronobserver`              |

### Optional Fields

These fields have defaults and are optional:

| Config Key             | Environment Variable   | Default | Description                  |
| ---------------------- | ---------------------- | ------- | ---------------------------- |
| `server.port`          | `SERVER_PORT`          | `8080`  | HTTP server port             |
| `server.read_timeout`  | `SERVER_READ_TIMEOUT`  | `15s`   | HTTP read timeout            |
| `server.write_timeout` | `SERVER_WRITE_TIMEOUT` | `15s`   | HTTP write timeout           |
| `database.timeout`     | `DATABASE_TIMEOUT`     | `10s`   | Database connection timeout  |
| `database.max_conns`   | `DATABASE_MAX_CONNS`   | `100`   | Maximum connection pool size |

## Usage Patterns

### Pattern 1: Using `.env` File (Development)

Create `.env` file in `backend/` directory:

```bash
DATABASE_URI=mongodb://localhost:27017
DATABASE_NAME=cronobserver
SERVER_PORT=8080
```

### Pattern 2: Using Environment Variables (Production)

```bash
export DATABASE_URI=mongodb://prod-server:27017
export DATABASE_NAME=cronobserver_prod
export SERVER_PORT=8080

go run cmd/server/main.go
```

### Pattern 3: Inline Environment Variables

```bash
DATABASE_URI=mongodb://localhost:27017 DATABASE_NAME=cronobserver go run cmd/server/main.go
```

### Pattern 4: Docker/Container Environment

```yaml
# docker-compose.yml
services:
  backend:
    environment:
      DATABASE_URI: mongodb://mongo:27017
      DATABASE_NAME: cronobserver
      SERVER_PORT: 8080
```

## Code Structure

### Configuration Package

The configuration system is organized in `internal/config/`:

```
internal/config/
├── config.go      # Type definitions
├── loader.go      # Viper loading logic
└── validation.go  # Validation and error handling
```

### Loading Configuration

Configuration is loaded at application startup in `cmd/server/main.go`:

```go
// Load configuration - fail fast if invalid
cfg, err := config.Load()
if err != nil {
    log.Fatalf("❌ Configuration Error:\n%v\n\nPlease check your environment variables or .env file.", err)
}
```

### Using Configuration

Configuration is injected into components:

```go
// Database connection
db, err := database.NewConnection(cfg.Database)

// HTTP server
srv := &http.Server{
    Addr:         ":" + cfg.Server.Port,
    ReadTimeout:  cfg.Server.ReadTimeout,
    WriteTimeout: cfg.Server.WriteTimeout,
}
```

## Fail-Fast Validation

The configuration system uses **fail-fast validation** - if required fields are missing, the application exits immediately with a clear error message:

```
❌ Missing Required Configuration

The following environment variables are required but not set:
  • DATABASE_URI
  • DATABASE_NAME

To fix this, either:
1. Set environment variables:
   export DATABASE_URI=your_value
   export DATABASE_NAME=your_value

2. Create a .env file in the backend directory:
   DATABASE_URI=your_value
   DATABASE_NAME=your_value

3. Or pass them when running:
   DATABASE_URI=your_value DATABASE_NAME=your_value go run cmd/server/main.go
```

## Adding New Configuration Values

To add a new configuration value:

### Step 1: Add to Config Struct

Edit `internal/config/config.go`:

```go
type ServerConfig struct {
    Port         string        `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    // Add new field
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}
```

### Step 2: Add Default (if optional) or Bind Env Var

Edit `internal/config/loader.go`:

```go
func setDefaults(v *viper.Viper) {
    // ... existing defaults

    // If optional, add default:
    v.SetDefault("server.idle_timeout", "60s")

    // Always bind environment variable:
    v.BindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")
}
```

### Step 3: Add Validation (if required)

Edit `internal/config/validation.go`:

```go
func (c *Config) Validate() error {
    // ... existing validation

    // If required:
    if c.Server.IdleTimeout == 0 {
        missing = append(missing, "SERVER_IDLE_TIMEOUT")
    }

    // If needs format validation:
    if c.Server.IdleTimeout < 0 {
        return fmt.Errorf("server.idle_timeout must be positive")
    }
}
```

### Step 4: Update Documentation

Update this file and `.env.example` with the new configuration option.

## Type Support

Viper automatically converts string values to the appropriate types:

- **Strings**: `"8080"` → `string`
- **Durations**: `"15s"`, `"5m"`, `"1h"` → `time.Duration`
- **Integers**: `"100"` → `int`
- **Booleans**: `"true"`, `"false"` → `bool`

## Best Practices

### 1. Required vs Optional

- **Required fields**: No default value, validated at startup
- **Optional fields**: Default value provided, can be overridden

### 2. Environment Variable Naming

Use uppercase with underscores, matching the config key:

- Config: `database.uri` → Env: `DATABASE_URI`
- Config: `server.read_timeout` → Env: `SERVER_READ_TIMEOUT`

### 3. Secrets Management

**Never commit secrets to `.env` files!**

- Use `.env.example` for documentation (no secrets)
- Add `.env` to `.gitignore`
- Use environment variables or secret managers in production

### 4. Development vs Production

- **Development**: Use `.env` file for convenience
- **Production**: Use environment variables or secret management systems (AWS Secrets Manager, HashiCorp Vault, etc.)

## Troubleshooting

### Issue: "Missing Required Configuration"

**Cause**: Required environment variables not set.

**Solution**:

1. Check `.env` file exists in `backend/` directory
2. Verify environment variables are set: `echo $DATABASE_URI`
3. Ensure variable names match exactly (case-sensitive)

### Issue: "Invalid duration format"

**Cause**: Duration string format incorrect.

**Solution**: Use valid Go duration format:

- `"10s"` ✅
- `"5m"` ✅
- `"1h30m"` ✅
- `"10"` ❌ (missing unit)

### Issue: Configuration not loading from `.env`

**Cause**: `.env` file not found or wrong location.

**Solution**:

1. Ensure `.env` is in `backend/` directory (same level as `go.mod`)
2. Check file permissions
3. Verify file format (no spaces around `=`)

### Issue: Environment variable not overriding `.env`

**Cause**: Variable name mismatch or not exported.

**Solution**:

1. Use `export VAR_NAME=value` (not just `VAR_NAME=value`)
2. Verify variable name matches exactly (case-sensitive)
3. Check variable is set: `env | grep DATABASE_URI`

## Example: Complete Setup

### `.env` File (Development)

```bash
# Database Configuration
DATABASE_URI=mongodb://localhost:27017
DATABASE_NAME=cronobserver
DATABASE_TIMEOUT=10s
DATABASE_MAX_CONNS=100

# Server Configuration
SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
```

### Production Environment Variables

```bash
export DATABASE_URI=mongodb+srv://user:pass@cluster.mongodb.net/dbname
export DATABASE_NAME=cronobserver_prod
export SERVER_PORT=8080
export SERVER_READ_TIMEOUT=30s
export SERVER_WRITE_TIMEOUT=30s
export DATABASE_TIMEOUT=15s
export DATABASE_MAX_CONNS=200
```

## References

- [Viper Documentation](https://github.com/spf13/viper)
- [Go Duration Format](https://pkg.go.dev/time#ParseDuration)
- [Environment Variables Best Practices](https://12factor.net/config)

