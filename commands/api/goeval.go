package api

import (
	"encoding/json"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var GoEvalSlash = discord.SlashCommandCreate{
	Name:        "eval_go",
	Description: "evaluate golang code",
}
var ExampleGoCode = `package main
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
	input, _ := e.Data.TextInputComponent("1-evalm")
	url1 := "https://play.golang.org/compile"

	payload := strings.NewReader(fmt.Sprintf("body=%s", url.QueryEscape(input.Value)))

	req, _ := http.NewRequest("POST", url1, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var bodym = map[string]any{}
	err := json.Unmarshal(body, &bodym)
	if err != nil {
		fmt.Println(err, string(body))
		return
	}
	if bodym["Errors"].(string) != "" {
		e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(OutputEmbed(input.Value, bodym["Errors"].(string))).Build())
		return
	}
	content := bodym["Events"].([]any)[0].(map[string]any)["Message"].(string)
	e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(OutputEmbed(input.Value, content)).Build())
}
func OutputEmbed(input, output string) discord.Embed {
	backticks := "```"
	builder := discord.NewEmbedBuilder().
		SetDescription(fmt.Sprintf(`
**INPUT**
%sgo
%s
%s
**OUTPUT**
%s
%s
%s
`, backticks, input, backticks, backticks, output, backticks))

	return builder.Build()

}
