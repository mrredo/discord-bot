package api

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var GoEvalSlash = discord.SlashCommandCreate{
	Name:        "eval_go",
	Description: "evaluate golang code",
}

func GoEval(e *events.ApplicationCommandInteractionCreate) {
	e.CreateModal(discord.NewModalCreateBuilder().SetTitle("Execute golang code").SetCustomID("1-goeval").Build())
}
