package jwt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// JWT interface
//go:generate mockery --name=JWT --filename=jwt_mock.go --inpackage
type JWT interface {
	// Generate a jwt token for the given admin user
	Generate(userID uint64) (string, error)

	// Parse the given JWT and return the admin user details
	Parse(token string) (uint64, error)
}

// jwtClaims struct
type jwtClaims struct {
	jwt.StandardClaims
}

// Config holds the JWT configuration
type Config struct {
	Expiry   uint64
	Key      string
	Issuer   string
	Audience string
}

var _ JWT = (*_jwt)(nil)

// _jwt implements token.Provider
type _jwt struct {
	config Config

	expiry time.Duration
	key    []byte
}

// JWT error declarations
var (
	ErrInvalidJWT = errors.New("invalid/expired token")
)

// New instantiates a new JWT instance
func New(config Config) *_jwt {
	return &_jwt{
		config: config,
		expiry: time.Duration(config.Expiry) * time.Second,
		key:    []byte(config.Key),
	}
}

// Generate a jwt token for the given user
func (t *_jwt) Generate(userID uint64) (string, error) {
	now := time.Now()
	expirationTime := now.Add(t.expiry)
	id := uuid.NewV4().String()
	claims := &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Id:        id,
			Subject:   fmt.Sprint(userID),
			Issuer:    t.config.Issuer,
			Audience:  t.config.Audience,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.key)
	if err != nil {
		return "", errors.Wrap(err, "unable to generate jwt token")
	}

	return tokenString, nil
}

func (t *_jwt) keyFunc(token *jwt.Token) (interface{}, error) {
	return t.key, nil
}

// Parse the given JWT and return the admin user details
func (t *_jwt) Parse(tokenString string) (uint64, error) {
	var userID uint64
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return userID, errors.Wrap(ErrInvalidJWT, "empty jwt token")
	}

	var claims jwtClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, t.keyFunc)
	if err != nil {
		return userID, errors.Wrap(err, "failed to parse token")
	}
	if !token.Valid || claims.Issuer != t.config.Issuer {
		return userID, errors.Wrap(ErrInvalidJWT, "invalid jwt token")
	}

	userID, err = strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return userID, errors.Wrap(err, "failed to parse userID")
	}

	return userID, nil
}
