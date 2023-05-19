package util

import "github.com/google/uuid"

func GenerateUUID() uuid.UUID {
	id := uuid.New()
	return id
}
