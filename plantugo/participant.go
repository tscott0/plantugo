package plantugo

import "github.com/fogleman/gg"

type Participant struct {
	Name          string
	Height, Width float64
	LineWidth     float64
	LineColour    string
	BgColour      string
	FontColour    string
	FontSize      float64
	FontAlign     gg.Align
}

func (p *Participant) Draw(s *SequenceDiagram, dc *gg.Context) {
	dc.Push()

	// Vertices
	dc.Push()
	dc.DrawLine(p.Width/2, p.Height, p.Width/2, p.Width/2+600)
	dc.SetLineWidth(3)
	dc.SetHexColor(p.LineColour)
	dc.SetDash(8, 8)
	dc.Stroke()
	dc.Pop()

	// Rectangle outline
	dc.Push()
	dc.DrawRectangle(0, 0, p.Width, p.Height)
	dc.SetLineWidth(p.LineWidth)
	dc.SetHexColor(p.LineColour)
	dc.StrokePreserve()
	// Rectangle background
	dc.SetHexColor(p.BgColour)
	dc.Fill()
	dc.Pop()

	// Text
	dc.SetHexColor(p.FontColour)
	if err := dc.LoadFontFace(s.Font, p.FontSize); err != nil {
		panic(err)
	}

	// A rough way to get the text to be centred vertically
	yPos := (p.Height / 2) - (p.FontSize / 3)

	dc.DrawStringWrapped(p.Name,
		0, yPos,
		0, 0, p.Width,
		1.5, p.FontAlign)

	dc.Pop()
}
