package bcrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_bcrypt_Init(t *testing.T) {
	tests := map[string]struct {
		cost        int
		expectedErr error
	}{
		"error case: value less than min cost": {
			cost:        bcrypt.MinCost - 1,
			expectedErr: ErrBcryptCostOutOfRange,
		},
		"error case: value greater than max cost": {
			cost:        bcrypt.MaxCost + 1,
			expectedErr: ErrBcryptCostOutOfRange,
		},
		"success case: value equal to min cost": {
			cost: bcrypt.MinCost,
		},
		"success case: value equal to max cost": {
			cost: bcrypt.MaxCost,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actualHasher, actualErr := New(Config{Cost: tc.cost})
			if tc.expectedErr == nil {
				assert.NoError(t, actualErr)
				assert.NotNil(t, actualHasher)
			} else {
				assert.EqualError(t, tc.expectedErr, actualErr.Error())
				assert.Nil(t, actualHasher)
			}
		})
	}
}

func Test_bcrypt_Hash(t *testing.T) {
	hasher, err := New(Config{Cost: bcrypt.DefaultCost})
	assert.NoError(t, err)
	assert.NotNil(t, hasher)

	hash, err := hasher.Hash("abcd")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}
func Test_bcrypt_Verify(t *testing.T) {
	hasher, err := New(Config{Cost: bcrypt.DefaultCost})
	assert.NoError(t, err)
	assert.NotNil(t, hasher)

	const str = "abcd"

	hash, err := hasher.Hash(str)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = hasher.Verify(hash, str)
	assert.NoError(t, err)
}
