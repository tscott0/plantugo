package plantugo

import "github.com/fogleman/gg"

type SequenceDiagram struct {
	OutputFile                                       string
	Height, Width                                    int
	MarginTop, MarginLeft, MarginRight, MarginBottom float64
	BgColour                                         string
	ParticipantTemplate                              Participant
	ParticipantSpacing                               float64
	Participants                                     map[string]Participant
	ParticipantKeys                                  []string
	Messages                                         []Message
	MessageSpacing                                   float64
}

func NewSequenceDiagram() SequenceDiagram {
	return SequenceDiagram{
		OutputFile:         "out.png",
		Width:              800,
		Height:             800,
		MarginTop:          50,
		MarginLeft:         50,
		MarginRight:        50,
		MarginBottom:       50,
		BgColour:           "#FFFFFF",
		ParticipantSpacing: 30,
		ParticipantTemplate: Participant{
			Height:     60,
			Width:      200,
			LineWidth:  5,
			LineColour: "#FFCC00",
			BgColour:   "#DDFFDD",
			Font:       "Roboto-Regular.ttf",
			FontColour: "#1100DD",
			FontSize:   28,
			FontAlign:  gg.AlignCenter,
		},
		Participants:    make(map[string]Participant),
		ParticipantKeys: make([]string, 0),
		Messages:        make([]Message, 0),
		MessageSpacing:  50,
	}
}

func (s *SequenceDiagram) AddParticipant(name string) {
	np := s.ParticipantTemplate
	s.Participants[name] = np
	s.ParticipantKeys = append(s.ParticipantKeys, name)
}

func (s *SequenceDiagram) Message(from, to, message string) {
	s.Messages = append(s.Messages, Message{from, to, message})
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
	dc := gg.NewContext(s.Width, s.Height)
	dc.SetHexColor(s.BgColour)
	dc.Clear()

	// Translate to provide top and left margins
	dc.Translate(s.MarginLeft, s.MarginTop)

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
		if err := dc.LoadFontFace(v.Font, v.FontSize); err != nil {
			panic(err)
		}

		// A rough way to get the text to be centred vertically
		yPos := (v.Height / 2) - (v.FontSize / 3)

		dc.DrawStringWrapped(name,
			0, yPos,
			0, 0, v.Width, 1.5, v.FontAlign)

		dc.Pop()

		dc.Translate(v.Width+s.ParticipantSpacing, 0)
	}
	dc.Pop()

	dc.Push()
	dc.Translate(0, s.ParticipantTemplate.Height+s.MessageSpacing)
	for _, m := range s.Messages {
		dc.Push()
		fromIndex := s.ParticipantIndex(m.From)
		toIndex := s.ParticipantIndex(m.To)

		dc.DrawLine((s.ParticipantTemplate.Width/2)+(float64(fromIndex)*(s.ParticipantTemplate.Width+s.ParticipantSpacing)),
			0,
			(s.ParticipantTemplate.Width/2)+(float64(toIndex)*(s.ParticipantTemplate.Width+s.ParticipantSpacing)),
			0)

		dc.SetLineWidth(2)
		dc.SetHexColor(s.ParticipantTemplate.LineColour)
		//dc.SetDash(8, 8)
		dc.Stroke()
		dc.Pop()

		dc.Translate(0, s.MessageSpacing)
	}
	dc.Pop()

	dc.SavePNG(s.OutputFile)
}
