// Package auth provides authentication and authorization support.
// Authentication: You are who you say you are.
// Authorization:  You have permission to do what you are requesting to do.
package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// ErrForbidden is returned when an auth issue is identified.
var ErrForbidden = errors.New("attempted action is not allowed")

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// Config represents information required to initialize auth.
type Config struct {
	Log    *logger.Logger
	DB     *sqlx.DB
	Secret string
	Issuer string
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	log     *logger.Logger
	secret  string
	userBus *userbus.Business
	issuer  string
}

// New creates an Auth instance to support authentication/authorization.
func New(cfg Config) (*Auth, error) {
	var userBus *userbus.Business
	if cfg.DB != nil {
		userBus = userbus.NewBusiness(cfg.Log, nil)
	}

	a := Auth{
		log:     cfg.Log,
		secret:  cfg.Secret,
		userBus: userBus,
		issuer:  cfg.Issuer,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return str, nil
}

// Authenticate validates the JWT token.
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return Claims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	tokenStr := bearerToken[7:]

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil || !token.Valid {
		return Claims{}, fmt.Errorf("authentication failed: %w", err)
	}

	if claims.Issuer != a.issuer {
		return Claims{}, errors.New("invalid token issuer")
	}

	return claims, nil
}

// Authorize checks if the user has the required role.
func (a *Auth) Authorize(ctx context.Context, claims Claims, userID uuid.UUID, requiredRole string) error {
	hasRole := false
	for _, role := range claims.Roles {
		if role == requiredRole {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return ErrForbidden
	}

	// Special case for admin-or-subject rule
	if requiredRole == "admin_or_subject" {
		if userID.String() != claims.Subject {
			return ErrForbidden
		}
	}

	return nil
}
