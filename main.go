package main

import (
	"context"
	"discord-bot/commands"
	"discord-bot/commands/api"
	"discord-bot/commands/games"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Info("disgo version: ", disgo.Version)

	client, err := disgo.New(os.Getenv("TOKEN"),
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(registerCommands),
		bot.WithEventListenerFunc(commandListener),
		bot.WithEventListenerFunc(api.HandleGoEvalInput),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	//if err = commands.RegisterCommands(client); err != nil {
	//	log.Fatal("error while registering commands: ", err)
	//}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
func registerCommands(event *events.MessageCreate) {
	if event.Message.Author.ID.String() == os.Getenv("OWNER_ID") && event.Message.Content == "!reg" {
		if err := commands.RegisterCommands(event.Client()); err != nil {
			log.Fatal("error while registering commands: ", err)
		}
	}
}
func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName() {
	case "eval_go":
		api.GoEval(event)

	case "nim":
		games.NimGame(event)
	default:
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Command not found").
			SetEphemeral(true).
			Build(),
		)
		if err != nil {
			event.Client().Logger().Error("error on sending response: ", err)
		}
	}

}
