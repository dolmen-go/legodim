package toypad_test

import (
	"image/color"

	"github.com/dolmen-go/legodim/toypad"
)

var _ color.Color = toypad.RGB{0, 0, 0}
