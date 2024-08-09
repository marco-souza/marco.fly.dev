package lua_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/lua"
	"github.com/stretchr/testify/assert"
)

func TestLua(t *testing.T) {
	l := lua.NewLuaService()

	t.Run("should run code if valid", func(t *testing.T) {
		snippet := "print(10^3)"
		output, err := l.Run(snippet)
		assert.Nil(t, err)
		assert.Contains(t, output, "1000")
	})

	t.Run("should error with invalid code", func(t *testing.T) {
		_, err := l.Run("(- (+ 1 1) 2)")
		assert.Error(t, err, "syntax error")

		_, err = l.Run("function js() { return 1+1 } js()")
		assert.Error(t, err, "syntax error")

		_, err = l.Run("invalid code")
		assert.Error(t, err, "syntax error")
	})
}
