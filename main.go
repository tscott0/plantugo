package main

import (
	"github.com/tscott0/plantugo/plantugo"
)

func main() {

	sd := plantugo.NewSequenceDiagram()
	sd.AddParticipant("Anna")
	sd.AddParticipant("Bob")
	sd.AddParticipant("Claire")

	sd.Message("Anna", "Bob", "Hello")
	sd.Message("Bob", "Anna", "How are you?")
	sd.Message("Claire", "Bob", "What time is it?")

	sd.Draw()
}
