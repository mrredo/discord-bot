package api

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"log"
)

var GoEvalSlash = discord.SlashCommandCreate{
	Name:        "eval_go",
	Description: "evaluate golang code",
}
var ExampleGoCode = `
package main
import "fmt"
func main() {
	fmt.Println("Hello World!")
}

`
var GoExampleModal = []discord.InteractiveComponent{
	discord.NewTextInput("1-evalm", discord.TextInputStyleParagraph, "Execute go code"),
}

func GoEval(e *events.ApplicationCommandInteractionCreate) {
	goe := GoExampleModal[0].(discord.TextInputComponent)
	goe.Value = ExampleGoCode
	GoExampleModal[0] = goe
	err := e.CreateModal(discord.NewModalCreateBuilder().AddActionRow(GoExampleModal...).SetTitle("Execute golang code").SetCustomID("1-goeval").Build())
	if err != nil {
		log.Println(err)
		return
	}
}
func HandleGoEvalInput(e *events.ModalSubmitInteractionCreate) {

}
