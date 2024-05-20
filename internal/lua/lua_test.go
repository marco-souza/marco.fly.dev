package lua

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLua(t *testing.T) {
	t.Run("should run code if valid", func(t *testing.T) {
		snippet := "return 1+1"
		output, err := Runtime.Run(snippet)
		assert.Nil(t, err)
		assert.Equal(t, output, "") // TODO: fix it later
	})

	t.Run("should error with invalid code", func(t *testing.T) {
		_, err := Runtime.Run("(- (+ 1 1) 2)")
		assert.Error(t, err, "syntax error")

		_, err = Runtime.Run("function js() { return 1+1 } js()")
		assert.Error(t, err, "syntax error")

		_, err = Runtime.Run("invalid code")
		assert.Error(t, err, "syntax error")
	})
}
