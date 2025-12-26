package misc

import (
	"image"
)

// Returns true if one rectangle is completely inside of another
func RectIsContained(inner, outer image.Rectangle) bool {
	return outer.Min.X <= inner.Min.X &&
		outer.Min.Y <= inner.Min.Y &&
		outer.Max.X >= inner.Max.X &&
		outer.Max.Y >= inner.Max.Y
}
