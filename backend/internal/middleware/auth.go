package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserInfo holds authenticated user information
type UserInfo struct {
	Email string
	Name  string
	Sub   string // User ID from JWT
}

// Context key for storing user info
const UserContextKey = "user"

// AuthMiddleware validates JWT tokens from NextAuth
// If the user's email is in the superAdmins list, they bypass JWT validation
func AuthMiddleware(jwtSecret string, superAdmins []string) gin.HandlerFunc {
	// Create a map for O(1) lookup
	superAdminMap := make(map[string]bool)
	for _, admin := range superAdmins {
		superAdminMap[strings.ToLower(strings.TrimSpace(admin))] = true
	}

	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// First, try to parse token without validation to extract email for super admin check
		parser := jwt.NewParser()
		unverifiedToken, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})

		var userEmail string
		var unverifiedClaims jwt.MapClaims
		if err == nil && unverifiedToken != nil {
			if claims, ok := unverifiedToken.Claims.(jwt.MapClaims); ok {
				unverifiedClaims = claims
				userEmail = getStringClaim(claims, "email")
				if userEmail == "" {
					if userObj, ok := claims["user"].(map[string]interface{}); ok {
						userEmail = getStringFromMap(userObj, "email")
					}
				}
			}
		}

		// Check if user is a super admin - if yes, bypass JWT validation
		if userEmail != "" && superAdminMap[strings.ToLower(strings.TrimSpace(userEmail))] {
			// Super admin - allow request without strict JWT validation
			userInfo := UserInfo{
				Email: userEmail,
			}
			if unverifiedClaims != nil {
				userInfo.Name = getStringClaim(unverifiedClaims, "name")
				userInfo.Sub = getStringClaim(unverifiedClaims, "sub")
				if userInfo.Name == "" {
					if userObj, ok := unverifiedClaims["user"].(map[string]interface{}); ok {
						userInfo.Name = getStringFromMap(userObj, "name")
					}
				}
			}
			c.Set(UserContextKey, userInfo)
			c.Next()
			return
		}

		// Not a super admin - proceed with normal JWT validation
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid or expired token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Extract user info from claims (NextAuth JWT format)
		userInfo := UserInfo{
			Email: getStringClaim(claims, "email"),
			Name:  getStringClaim(claims, "name"),
			Sub:   getStringClaim(claims, "sub"),
		}

		// If email is missing, try to get it from user object in token
		if userInfo.Email == "" {
			if userObj, ok := claims["user"].(map[string]interface{}); ok {
				userInfo.Email = getStringFromMap(userObj, "email")
				userInfo.Name = getStringFromMap(userObj, "name")
			}
		}

		// Store user info in context for handlers to access
		c.Set(UserContextKey, userInfo)

		// Continue to next handler
		c.Next()
	}
}

// GetUserFromContext extracts user info from gin context
func GetUserFromContext(c *gin.Context) (*UserInfo, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}

	userInfo, ok := user.(UserInfo)
	if !ok {
		return nil, false
	}

	return &userInfo, true
}

// Helper to safely extract string claims
func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Helper to safely extract string from map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
