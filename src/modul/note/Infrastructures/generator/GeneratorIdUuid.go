package generator

import (
	"fmt"

	"github.com/google/uuid"
)

type GeneratorIdUuid struct{}

func (h *GeneratorIdUuid) Generator() string {
	id := fmt.Sprintf("note-%s", uuid.New().String())

	return id
}
