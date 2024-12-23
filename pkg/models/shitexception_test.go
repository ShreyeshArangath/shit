package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShitException(t *testing.T) {
	exception := &ShitException{Message: "test"}
	// Check if exception conforms to the error interface
	var err error = exception
	assert.Implements(t, (*error)(nil), err)
	assert.Equal(t, "ShitException: test", exception.Error())
}
