package plantugo

import "github.com/fogleman/gg"

type Participant struct {
	Height, Width float64
	LineWidth     float64
	LineColour    string
	BgColour      string
	FontColour    string
	FontSize      float64
	FontAlign     gg.Align
}
