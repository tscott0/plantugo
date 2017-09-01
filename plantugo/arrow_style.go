package plantugo

type ArrowStyle uint8

const (
	SolidArrow ArrowStyle = iota
	EmptyArrow
)

func (as ArrowStyle) StyleMessage(m *Message) {
	m.ArrowStyle = as
}

func (as ArrowStyle) Width() float64 {

	return 1
}
