package di_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	type Test struct {
		name string
	}

	// Add your test here
	t.Run("Register an structure", func(t *testing.T) {
		di.Injectable(Test{name: "test"})
	})

	// inject an instance of Test
	t.Run("Get an structure", func(t *testing.T) {
		t1 := di.Inject(Test{})
		t2 := di.Inject(Test{})

		assert.Equal(t, t1.name, t2.name)
		assert.Equal(t, t1, t2)
	})

	// clean the container
	t.Run("Clean the container", func(t *testing.T) {
		di.Clean()
		assert.Empty(t, di.Inject(Test{}))
	})
}
