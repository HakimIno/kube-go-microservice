package middleware

import (
	"context"
	"strings"
	"time"

	"kube/pkg/errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(secretKey string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "Authorization header required", "Please provide a valid authorization token"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "Bearer token required", "Please provide a valid Bearer token"))
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(errors.ErrCodeUnauthorized, "Invalid signing method", "Token signing method is invalid")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "Invalid token", "Token validation failed"))
			c.Abort()
			return
		}

		if !token.Valid {
			errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "Invalid token", "Token is not valid"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "Invalid token claims", "Token claims are invalid"))
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)

		c.Next(ctx)
	}
}

// OptionalAuthMiddleware allows requests to proceed even without valid auth
func OptionalAuthMiddleware(secretKey string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.Next(ctx)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next(ctx)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(errors.ErrCodeUnauthorized, "Invalid signing method", "Token signing method is invalid")
			}
			return []byte(secretKey), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*Claims); ok {
				c.Set("user_id", claims.UserID)
				c.Set("email", claims.Email)
				c.Set("role", claims.Role)
				c.Set("token", tokenString)
			}
		}

		c.Next(ctx)
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRole string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userRole, exists := c.Get("role")
		if !exists {
			errors.SendError(c, errors.New(errors.ErrCodeForbidden, "Role required", "User role not found in context"))
			c.Abort()
			return
		}

		if userRole != requiredRole {
			errors.SendError(c, errors.New(errors.ErrCodeForbidden, "Insufficient permissions", "You don't have permission to access this resource"))
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// GenerateToken creates a new JWT token
func GenerateToken(userID uint, email, role, secretKey string, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// RefreshToken creates a new token with extended expiration
func RefreshToken(tokenString, secretKey string, expiration time.Duration) (string, error) {
	// Parse existing token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New(errors.ErrCodeUnauthorized, "Invalid token", "Cannot refresh invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New(errors.ErrCodeUnauthorized, "Invalid token claims", "Token claims are invalid")
	}

	// Create new token with extended expiration
	return GenerateToken(claims.UserID, claims.Email, claims.Role, secretKey, expiration)
}
