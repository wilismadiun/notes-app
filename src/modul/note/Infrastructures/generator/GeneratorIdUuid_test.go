package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	generateHandler := GeneratorIdUuid{}

	t.Run("should generate uuid with note prefix", func(t *testing.T) {
		id := generateHandler.Generator()

		assert.True(t, strings.HasPrefix(id, "note-"))
	})
}
