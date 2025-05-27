package gokick_test

import (
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	t.Run("without description", func(t *testing.T) {
		err := gokick.NewError(401, "not authorized")

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, 401, kickError.Code())
		assert.Equal(t, "not authorized", kickError.Message())
		assert.Empty(t, kickError.Description())
		assert.EqualError(t, err, "Error 401: not authorized")
	})

	t.Run("with description", func(t *testing.T) {
		err := gokick.NewError(401, "not authorized").WithDescription("invalid scope")

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, 401, kickError.Code())
		assert.Equal(t, "not authorized", kickError.Message())
		assert.Equal(t, "invalid scope", kickError.Description())
		assert.EqualError(t, err, "Error 401: not authorized (invalid scope)")
	})
}
