package engine

import "image/color"

// RGBA constructs a color.RGBA with bytes
func RGBA(r, g, b, a uint8) color.RGBA { return color.RGBA{r, g, b, a} }