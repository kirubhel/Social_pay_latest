package entity

import (
	"image/color"

	"github.com/google/uuid"
)

type Tag struct {
	Id    uuid.UUID    `json:"id"`
	Name  string       `json:"name"`
	Color color.RGBA64 `json:"color"`
}
