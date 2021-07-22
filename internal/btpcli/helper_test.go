package btpcli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstElementOrDefault(t *testing.T) {
	t.Parallel()
	t.Run("returns default, slice is empty", func(t *testing.T) {
		val := firstElementOrDefault([]int{}, 1)

		assert.Equal(t, 1, val)
	})
	t.Run("return first element", func(t *testing.T) {
		val := firstElementOrDefault([]int{3, 2, 1}, 1)

		assert.Equal(t, 3, val)
	})
}

func TestNthElementOrDefault(t *testing.T) {
	t.Parallel()
	t.Run("returns default, slice is empty", func(t *testing.T) {
		val := nthElementOrDefault([]int{}, 0, 1)

		assert.Equal(t, 1, val)
	})
	t.Run("return default, slice has not enough elements", func(t *testing.T) {
		val := nthElementOrDefault([]int{3}, 2, 1)

		assert.Equal(t, 1, val)
	})
	t.Run("return nth element", func(t *testing.T) {
		val := nthElementOrDefault([]int{3, 2, 1}, 1, 1)

		assert.Equal(t, 2, val)
	})
}
