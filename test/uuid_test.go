package test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestUUID(t *testing.T) {
	uuid := uuid.New()
	fmt.Println(uuid.String())
}
