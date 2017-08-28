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

func (s *SequenceDiagram) Participant(name string) {
	np := s.ParticipantTemplate
	np.Name = name
	s.Participants[name] = np

	s.ParticipantKeys = append(s.ParticipantKeys, name)
	fmt.Printf("Added %q. Total participants: %v\n", name, len(s.Participants))
}

func (s *SequenceDiagram) Message(from, to, message string, styles ...MessageStyler) Message {
	m := Message{
		From:       from,
		To:         to,
		Message:    message,
		LineStyle:  SolidLine,
		ArrowStyle: SolidArrow,
	}

	for _, style := range styles {
		style.StyleMessage(&m)
	}

	s.Messages = append(s.Messages, m)

	// Insert participants if not explicitly created
	s.ParticipantIndex(from)
	s.ParticipantIndex(to)

	return m
}

func (s *SequenceDiagram) MessageSelf(participant, message string, styles ...MessageStyler) Message {
	m := Message{
		From:       participant,
		To:         participant,
		Message:    message,
		LineStyle:  SolidLine,
		ArrowStyle: SolidArrow,
	}

	for _, style := range styles {
		style.StyleMessage(&m)
	}

	s.Messages = append(s.Messages, m)

	// Insert participants if not explicitly created
	s.ParticipantIndex(participant)

	return m
}

func (s *SequenceDiagram) ParticipantIndex(name string) int {
	if _, ok := s.Participants[name]; !ok {
		s.Participant(name)
	}

	for i, v := range s.ParticipantKeys {
		if v == name {
			return i
		}
	}

	return -1
}

func (s *SequenceDiagram) Dimensions() (int, int) {
	// Calculate the width of the image
	w := int(2 * s.Margin)
	w += len(s.Participants) * int(s.ParticipantTemplate.Width)
	w += (len(s.Participants) - 1) * int(s.ParticipantSpacing)

	w += int(s.ParticipantTemplate.Width) // TODO: Temporary extra width

	// TODO: Calculate the height of the image
	h := 800

	return w, h
}

func (s *SequenceDiagram) Draw() {
	dc := gg.NewContext(s.Dimensions())
	dc.SetHexColor(s.BgColour)
	dc.Clear()

	// Translate to provide top and left margins
	dc.Translate(s.Margin, s.Margin)

	// PARTICIPANTS
	dc.Push()
	for _, name := range s.ParticipantKeys {
		p := s.Participants[name]
		p.Draw(s, dc)
		dc.Translate(p.Width+s.ParticipantSpacing, 0)
	}
	dc.Pop()

	// MESSAGES
	dc.Push()
	dc.Translate(0, s.ParticipantTemplate.Height+s.MessageSpacing)
	for _, m := range s.Messages {
		m.Draw(s, dc)
	}
	dc.Pop()

	dc.SavePNG(s.OutputFile)
}
