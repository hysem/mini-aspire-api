package jwt_test

import (
	"testing"
	"time"

	"github.com/hysem/mini-aspire-api/app/core/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWT_Generate(t *testing.T) {
	j := jwt.New(jwt.Config{
		Expiry:   2,
		Key:      "my_secret_key",
		Issuer:   "test.local",
		Audience: "test.local",
	})

	jwtString, err := j.Generate(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtString)
}

func TestJWT_Parse(t *testing.T) {
	j := jwt.New(jwt.Config{
		Expiry:   2,
		Key:      "my_secret_key",
		Issuer:   "test.local",
		Audience: "test.local",
	})

	expectedUserID := uint64(10)

	jwtString, err := j.Generate(expectedUserID)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtString)

	t.Run("with a fresh jwt", func(t *testing.T) {
		time.Sleep(1 * time.Second)
		actualUserID, err := j.Parse(jwtString)
		assert.NoError(t, err)
		assert.Equal(t, expectedUserID, actualUserID)
	})

	t.Run("with an expired jwt", func(t *testing.T) {
		time.Sleep(2 * time.Second)
		actualUserID, err := j.Parse(jwtString)
		assert.Contains(t, err.Error(), `token is expired by 1s`)
		assert.Zero(t, actualUserID)
	})

	t.Run("with a token generated by another issuer", func(t *testing.T) {
		anotherJ := jwt.New(jwt.Config{
			Expiry:   2,
			Key:      "my_secret_key",
			Issuer:   "test1.local",
			Audience: "test1.local",
		})
		anotherJWTstring, err := anotherJ.Generate(expectedUserID)
		assert.NoError(t, err)
		assert.NotEmpty(t, anotherJWTstring)

		actualUserID, err := j.Parse(jwtString)
		assert.Contains(t, err.Error(), `failed to parse token`)
		assert.Zero(t, actualUserID)
	})

}
