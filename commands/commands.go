package commands

import (
	"discord-bot/commands/api"
	"discord-bot/commands/games"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

var Commands = []discord.ApplicationCommandCreate{}

// add dictionary support so auto detects
func NewCommand(command discord.SlashCommandCreate) {
	Commands = append(Commands, command)
}
func RegisterCommands(c bot.Client) error {
	LoadCommands()
	_, err := c.Rest().SetGlobalCommands(c.ApplicationID(), Commands)
	return err
}
func LoadCommands() {
	NewCommand(games.NimSlashCommand)
	NewCommand(api.GoEvalSlash)
}
