package bcrypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt interface
//go:generate mockery --name=Bcrypt --filename=bcrypt_mock.go --inpackage
type Bcrypt interface {
	// Hash generates bcrypt hash from password
	Hash(str string) (string, error)

	// Verify bcrypt hash and string
	Verify(hash string, str string) error
}

// Bcrypt error declarations
var (
	ErrBcryptCostOutOfRange = fmt.Errorf("cost should be within [%d, %d]", bcrypt.MinCost, bcrypt.MaxCost)
)

// Config holds the bcrypt configuration
type Config struct {
	Cost int
}

var _ Bcrypt = (*_bcrypt)(nil)

// _bcrypt implements Bcrypt interface
type _bcrypt struct {
	config Config
}

// New instantitates a _bcrypt instance
func New(config Config) (*_bcrypt, error) {
	if config.Cost > bcrypt.MaxCost || config.Cost < bcrypt.MinCost {
		return nil, ErrBcryptCostOutOfRange
	}
	return &_bcrypt{
		config: config,
	}, nil
}

// Verify bcrypt hash and string
func (h *_bcrypt) Verify(hash string, str string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}

// Hash generates bcrypt hash from password
func (h *_bcrypt) Hash(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), h.config.Cost)
	return string(b), err
}
