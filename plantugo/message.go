package plantugo

import (
	"github.com/fogleman/gg"
)

type Message struct {
	From, To, Message string
	LineStyle         LineStyle
	ArrowStyle        ArrowStyle
}

type MessageStyler interface {
	StyleMessage(m *Message)
}

func (m *Message) SetStyle(s LineStyle) {
	m.LineStyle = s
}

func (m *Message) Draw(s *SequenceDiagram, dc *gg.Context) {
	if m.From == m.To {
		m.drawSelf(s, dc)
	} else {
		m.draw(s, dc)
	}
}

func (m *Message) draw(s *SequenceDiagram, dc *gg.Context) {
	// Message line
	fromIndex := s.ParticipantIndex(m.From)
	toIndex := s.ParticipantIndex(m.To)

	originX := s.ParticipantTemplate.Width / 2
	originX += float64(fromIndex) * (s.ParticipantTemplate.Width + s.ParticipantSpacing)

	destinationX := s.ParticipantTemplate.Width / 2
	destinationX += float64(toIndex) * (s.ParticipantTemplate.Width + s.ParticipantSpacing)

	leftToRight := false
	if originX < destinationX {
		leftToRight = true
	}

	// Adjust for the width of the vertical activity line
	lineWidth := float64(4)
	if leftToRight {
		originX += lineWidth / 2
		destinationX -= lineWidth / 2
	} else {
		originX -= lineWidth / 2
		destinationX += lineWidth / 2
	}

	dc.DrawLine(originX,
		0,
		destinationX,
		0)

	dc.SetLineWidth(s.MessageLineWidth * 10)
	dc.SetHexColor(s.MessageLineColour)
	dc.SetLineCapButt()

	if m.LineStyle == DashedLine {
		dc.SetDash(4, 4)
	} else {
		dc.SetDash()
	}

	dc.Stroke()

	// Arrow
	dc.Push()
	dc.Translate(destinationX, 0)

	dc.LineTo(0, 0)
	if leftToRight {
		dc.LineTo(-16, -8)
		dc.LineTo(-16, 8)
	} else {
		dc.LineTo(16, -8)
		dc.LineTo(16, 8)
	}
	dc.LineTo(0, 0)

	dc.SetLineWidth(s.MessageLineWidth)

	switch m.ArrowStyle {
	case EmptyArrow:
		// TODO: Fix arrow styling for empty arrow
		dc.SetHexColor(s.ParticipantTemplate.LineColour)
	default:
		dc.SetHexColor(s.MessageLineColour)
	}
	dc.Fill()

	dc.SetHexColor(s.MessageLineColour)
	dc.StrokePreserve()
	dc.Pop()

	// Message text
	dc.SetHexColor(s.MessageFontColour)
	if err := dc.LoadFontFace(s.Font, s.MessageFontSize); err != nil {
		panic(err)
	}

	// Attempt to align font vertically above the message line
	yMessageOffset := -10 - s.MessageFontSize/2 - s.MessageLineWidth

	if originX < destinationX {
		dc.DrawStringWrapped(m.Message,
			originX, yMessageOffset,
			-0.05, 0, s.ParticipantTemplate.Width,
			1.5, gg.AlignLeft)
	} else {
		dc.DrawStringWrapped(m.Message,
			originX, yMessageOffset,
			1.05, 0, s.ParticipantTemplate.Width,
			1.5, gg.AlignRight)
	}

	// Move down
	dc.Translate(0, s.MessageSpacing)
}

func (m *Message) drawSelf(s *SequenceDiagram, dc *gg.Context) {
	// Message line
	fromIndex := s.ParticipantIndex(m.From)

	originX := s.ParticipantTemplate.Width / 2
	originX += float64(fromIndex) * (s.ParticipantTemplate.Width + s.ParticipantSpacing)

	// Adjust for the width of the vertical activity line
	lineWidth := float64(4)
	originX += lineWidth / 2

	dc.Push()
	dc.Translate(originX, 0)

	loopWidth := s.ParticipantTemplate.Width / 3
	loopHeight := s.ParticipantTemplate.Width / 6 // Also determines vertical spacing

	dc.LineTo(0, 0)
	dc.LineTo(loopWidth, 0)
	dc.LineTo(loopWidth, loopHeight)
	dc.LineTo(0, loopHeight)

	dc.SetLineWidth(s.MessageLineWidth)
	dc.SetHexColor(s.MessageLineColour)

	switch m.LineStyle {
	case DashedLine:
		dc.SetDash(4, 4)
	default:
		dc.SetDash()
	}

	dc.Stroke()

	// Arrow
	dc.Translate(0, loopHeight)

	dc.LineTo(0, 0)
	dc.LineTo(16, -8)
	dc.LineTo(16, 8)
	dc.LineTo(0, 0)

	dc.SetLineWidth(s.MessageLineWidth)
	dc.SetHexColor(s.MessageLineColour)

	switch m.ArrowStyle {
	case EmptyArrow:
		// Do nothing
	default:
		dc.Fill()
	}

	dc.StrokePreserve()
	dc.Pop()

	// Message text
	dc.SetHexColor(s.MessageFontColour)
	if err := dc.LoadFontFace(s.Font, s.MessageFontSize); err != nil {
		panic(err)
	}

	// Attempt to align font vertically above the message line
	yMessageOffset := -10 - s.MessageFontSize/2 - s.MessageLineWidth

	dc.DrawStringWrapped(m.Message,
		originX, yMessageOffset,
		-0.05, 0, s.ParticipantTemplate.Width,
		1.5, gg.AlignLeft)

	// Move down
	dc.Translate(0, s.MessageSpacing+loopHeight)
}
