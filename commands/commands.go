package commands

import (
	"discord-bot/commands/games"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

var Commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "bozo",
		Description: "you big bozo",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "who",
				Description: "What to say",
				Required:    true,
			},
		},
	},
}

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
}
