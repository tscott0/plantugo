package main

import (
	"github.com/tscott0/plantugo/plantugo"
)

func main() {

	sd := plantugo.NewSequenceDiagram()
	sd.Participant("Anna")
	sd.Participant("Bob")

	sd.Message("Anna", "Bob", "Hello")
	sd.Message("Bob", "Anna", "How are you?",
		plantugo.DashedLine, plantugo.EmptyArrow)
	sd.Message("Anna", "Bob", "Fine thanks")
	sd.Message("Claire", "Bob", "What time is it?")

	sd.MessageSelf("Claire", "Making sandwich",
		plantugo.DashedLine)
	sd.Message("Bob", "Bob", "I don't know")

	sd.Message("Claire", "Anna", "All the way back", plantugo.EmptyArrow)

	sd.Draw()
}
