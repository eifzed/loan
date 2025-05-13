package util

import (
	"github.com/google/uuid"
)

var (
	u = uuid.New()
)

func GenerateUUID() string {
	return u.String()
}
