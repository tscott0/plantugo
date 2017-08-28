package plantugo

type LineStyle uint8

const (
	SolidLine LineStyle = iota
	DashedLine
)

func (ls LineStyle) StyleMessage(m *Message) {
	m.LineStyle = ls
}
