package lua_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/lua"
	"github.com/stretchr/testify/assert"
)

func TestLua(t *testing.T) {
	t.Run("should run code if valid", func(t *testing.T) {
		snippet := "print(10^3)"
		output, err := lua.Run(snippet)
		assert.Nil(t, err)
		assert.Contains(t, output, "1000")
	})

	t.Run("should error with invalid code", func(t *testing.T) {
		_, err := lua.Run("(- (+ 1 1) 2)")
		assert.Error(t, err, "syntax error")

		_, err = lua.Run("function js() { return 1+1 } js()")
		assert.Error(t, err, "syntax error")

		_, err = lua.Run("invalid code")
		assert.Error(t, err, "syntax error")
	})
}
