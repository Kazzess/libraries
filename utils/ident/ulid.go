package ident

import (
	"github.com/oklog/ulid/v2"
)

type ULIDGenerator struct{}

func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{}
}

func (g *ULIDGenerator) GenerateID() string {
	return ulid.Make().String()
}
