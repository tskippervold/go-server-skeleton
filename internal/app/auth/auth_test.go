package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashAndCompare(t *testing.T) {
	assert := assert.New(t)

	password := "1234_password"

	hashedPassword, err := hashPassword(password)
	assert.NoError(err)

	err = compareHashAndPassword(hashedPassword, password)
	assert.NoError(err)
}
