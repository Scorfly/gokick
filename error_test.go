package gokick_test

import (
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	err := gokick.NewError(401, "not authorized")

	var kickError gokick.Error
	require.ErrorAs(t, err, &kickError)
	assert.Equal(t, 401, kickError.Code())
	assert.Equal(t, "not authorized", kickError.Message())
	assert.EqualError(t, err, "Error 401: not authorized")
}
