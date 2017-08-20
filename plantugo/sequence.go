package plantugo

import (
	"fmt"
	"github.com/fogleman/gg"
)

type SequenceDiagram struct {
	OutputFile          string
	Margin              float64
	BgColour            string
	Font                string
	ParticipantTemplate Participant
	ParticipantSpacing  float64
	Participants        map[string]Participant
	ParticipantKeys     []string
	Messages            []Message
	MessageSpacing      float64
	MessageLineColour   string
	MessageLineWidth    float64
	MessageFontColour   string
	MessageFontSize     float64
}

func NewSequenceDiagram() SequenceDiagram {
	return SequenceDiagram{
		OutputFile:         "out.png",
		Margin:             50,
		BgColour:           "#FFFFFF",
		Font:               "Roboto-Regular.ttf",
		ParticipantSpacing: 30,
		ParticipantTemplate: Participant{
			Height:     60,
			Width:      200,
			LineWidth:  5,
			LineColour: "#FFCC00",
			BgColour:   "#DDFFDD",
			FontColour: "#1100DD",
			FontSize:   28,
			FontAlign:  gg.AlignCenter,
		},
		Participants:      make(map[string]Participant),
		ParticipantKeys:   make([]string, 0),
		Messages:          make([]Message, 0),
		MessageSpacing:    50,
		MessageLineColour: "#119911",
		MessageLineWidth:  2,
		MessageFontColour: "#771111",
		MessageFontSize:   18,
	}
}

func (s *SequenceDiagram) AddParticipant(name string) {
	np := s.ParticipantTemplate
	s.Participants[name] = np
	s.ParticipantKeys = append(s.ParticipantKeys, name)
	fmt.Printf("Added %q. Total participants: %v\n", name, len(s.Participants))
}

func (s *SequenceDiagram) Message(from, to, message string) {
	s.Messages = append(s.Messages, Message{from, to, message})

	// Insert participants if not explicitly created
	s.ParticipantIndex(from)
	s.ParticipantIndex(to)
}

func (s *SequenceDiagram) ParticipantIndex(name string) int {
	if _, ok := s.Participants[name]; !ok {
		s.AddParticipant(name)
	}

	for i, v := range s.ParticipantKeys {
		if v == name {
			return i
		}
	}

	return -1
}

func (s *SequenceDiagram) Draw() {
	// Calculate the width of the image
	w := int(2 * s.Margin)
	w += len(s.Participants) * int(s.ParticipantTemplate.Width)
	w += (len(s.Participants) - 1) * int(s.ParticipantSpacing)

	// TODO: Calculate the height of the image
	h := 800

	dc := gg.NewContext(w, h)
	dc.SetHexColor(s.BgColour)
	dc.Clear()

	// Translate to provide top and left margins
	dc.Translate(s.Margin, s.Margin)

	dc.Push()
	for _, name := range s.ParticipantKeys {
		dc.Push()

		v := s.Participants[name]

		// Vertices
		dc.Push()
		dc.DrawLine(v.Width/2, v.Height, v.Width/2, v.Width/2+600)
		dc.SetLineWidth(3)
		dc.SetHexColor(v.LineColour)
		dc.SetDash(8, 8)
		dc.Stroke()
		dc.Pop()

		// Rectangle outline
		dc.Push()
		dc.DrawRectangle(0, 0, v.Width, v.Height)
		dc.SetLineWidth(v.LineWidth)
		dc.SetHexColor(v.LineColour)
		dc.StrokePreserve()
		// Rectangle background
		dc.SetHexColor(v.BgColour)
		dc.Fill()
		dc.Pop()

		// Text
		dc.SetHexColor(v.FontColour)
		if err := dc.LoadFontFace(s.Font, v.FontSize); err != nil {
			panic(err)
		}

		// A rough way to get the text to be centred vertically
		yPos := (v.Height / 2) - (v.FontSize / 3)

		dc.DrawStringWrapped(name,
			0, yPos,
			0, 0, v.Width,
			1.5, v.FontAlign)

		dc.Pop()

		dc.Translate(v.Width+s.ParticipantSpacing, 0)
	}
	dc.Pop()

	dc.Push()
	dc.Translate(0, s.ParticipantTemplate.Height+s.MessageSpacing)
	for _, m := range s.Messages {
		dc.Push()

		// Message line
		fromIndex := s.ParticipantIndex(m.From)
		toIndex := s.ParticipantIndex(m.To)

		originX := s.ParticipantTemplate.Width / 2
		originX += float64(fromIndex) * (s.ParticipantTemplate.Width + s.ParticipantSpacing)

		destinationX := s.ParticipantTemplate.Width / 2
		destinationX += float64(toIndex) * (s.ParticipantTemplate.Width + s.ParticipantSpacing)

		dc.DrawLine(originX,
			0,
			destinationX,
			0)

		dc.SetLineWidth(s.MessageLineWidth)
		dc.SetHexColor(s.MessageLineColour)
		//dc.SetDash(8, 8)

		dc.Stroke()

		// Arrow
		dc.Push()
		dc.Translate(destinationX, 0)

		dc.LineTo(0, 0)
		if originX < destinationX {
			dc.LineTo(-16, -8)
			dc.LineTo(-16, 8)
		} else {
			dc.LineTo(16, -8)
			dc.LineTo(16, 8)
		}
		dc.LineTo(0, 0)

		dc.SetLineWidth(s.MessageLineWidth)
		dc.SetHexColor(s.MessageLineColour)
		dc.Fill()
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

		dc.Pop()

		dc.Translate(0, s.MessageSpacing)
	}
	dc.Pop()

	dc.SavePNG(s.OutputFile)
}
