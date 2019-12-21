package utils

import (
	"math/rand"
	"testing"
	"time"
)

func TestGenerateId(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Log(GenerateId(5))
}
