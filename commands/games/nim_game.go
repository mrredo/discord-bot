package games

import (
	"context"
	"discord-bot/functions"
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"log"
	"strconv"
	"strings"
	"time"
)

func newNimEmbed(stones int, user1, user2 snowflake.ID, player1Turn bool, won ...bool) discord.Embed {

	var turn string
	if player1Turn {
		turn = fmt.Sprintf("<@%s>", user1.String())
	} else {
		turn = fmt.Sprintf("<@%s>", user2.String())
	}
	if len(won) != 0 {
		winnerid := ""
		if !won[0] {
			winnerid = user1.String()
		} else {
			winnerid = user2.String()
		}
		return discord.NewEmbedBuilder().
			SetTitle("Nim Game").
			SetDescriptionf("winner: <@%s> \n<@%s> vs <@%s>\n stones: %d\nturn: %s", winnerid, user1.String(), user2.String(), stones, turn).
			SetColor(8552316).
			Build()
	}

	return discord.NewEmbedBuilder().
		SetTitle("Nim Game").
		SetDescriptionf("<@%s> vs <@%s>\n stones: %d\nturn: %s", user1.String(), user2.String(), stones, turn).
		SetColor(8552316).
		//SetField(0, "Stones left", fmt.Sprintf("%d stones left", stones), true).
		Build()
}

var NimButtons = []discord.InteractiveComponent{
	discord.NewPrimaryButton("-1", "1-nim"),
	discord.NewPrimaryButton("-2", "2-nim"),
	discord.NewPrimaryButton("-3", "3-nim"),
	discord.NewDangerButton("resign", "resign-nim"),
}

func NimGame(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()
	user1, user2 := e.User(), data.User("user")
	randomUserNum := functions.RandomInRange(1, 2)
	if randomUserNum == 1 {
		user2, user1 = e.User(), data.User("user")
	}
	if user1.ID == user2.ID {
		e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("You can't duel yourself").SetEphemeral(true).Build())
	} else if user2.Bot {
		e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("I don't think a bot will play").SetEphemeral(true).Build())
	}
	var stones = functions.RandomInRange(20, 40)
	if err := e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(newNimEmbed(stones, user1.ID, user2.ID, true)).AddActionRow(NimButtons...).Build()); err != nil {
		log.Println(err)
	}
	isPlayer1Turn := true
	go func() {

		ch, cls := bot.NewEventCollector(e.Client(), func(e1 *events.ComponentInteractionCreate) bool {
			return e1.User().ID == user1.ID || e1.User().ID == user2.ID
		})
		defer cls()
		ctx, clsCtx := context.WithTimeout(context.Background(), 3*time.Minute)
		defer clsCtx()
		for {
			select {
			case <-ctx.Done(): //cancelled = resigns
				disabledButtons := []discord.InteractiveComponent{}
				for _, v := range NimButtons {

					disabledButtons = append(disabledButtons, v.(discord.ButtonComponent).AsDisabled())
				}
				e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateBuilder().ClearContainerComponents().AddActionRow(disabledButtons...).Build())
				return

			case bEvent := <-ch:
				if bEvent.ButtonInteractionData().CustomID() == "resign-nim" {
					bEvent.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Successfully resigned.").SetEphemeral(true).Build())
					disabledButtons := []discord.InteractiveComponent{}
					for _, v := range NimButtons {

						disabledButtons = append(disabledButtons, v.(discord.ButtonComponent).AsDisabled())
					}
					isPlayer1Turn = !isPlayer1Turn
					e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateBuilder().ClearContainerComponents().AddActionRow(disabledButtons...).SetEmbeds(newNimEmbed(stones, user1.ID, user2.ID, isPlayer1Turn, !isPlayer1Turn)).Build())
					return
				}
				if bEvent.User().ID != user1.ID && bEvent.User().ID != user2.ID {
					err := bEvent.CreateMessage(discord.NewMessageCreateBuilder().SetContent("This is not your game, create your own.").SetEphemeral(true).Build())
					if err != nil {
						log.Println(err)
					}
					break
				}
				if isPlayer1Turn && bEvent.User().ID != user1.ID || !isPlayer1Turn && bEvent.User().ID != user2.ID {
					err := bEvent.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Its not your turn, wait.").SetEphemeral(true).Build())
					if err != nil {
						log.Println(err)
					}
					break
				} /*else if !isPlayer1Turn && bEvent.User().ID != user2.ID {
					return
				}*/

				num, err := strconv.Atoi(strings.Split(bEvent.ButtonInteractionData().CustomID(), "-")[0])
				if err != nil {
					panic(err)
				}
				stones -= num

				bEvent.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Successfully removed %d stones.", num).SetEphemeral(true).Build())
				isPlayer1Turn = !isPlayer1Turn
				if stones <= 0 {
					stones = 0
					disabledButtons := []discord.InteractiveComponent{}
					for _, v := range NimButtons {

						disabledButtons = append(disabledButtons, v.(discord.ButtonComponent).AsDisabled())
					}
					e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateBuilder().ClearContainerComponents().AddActionRow(disabledButtons...).SetEmbeds(newNimEmbed(stones, user1.ID, user2.ID, isPlayer1Turn, isPlayer1Turn)).Build())
					return
				} else {
					e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(newNimEmbed(stones, user1.ID, user2.ID, isPlayer1Turn)).Build())
				}

			}
		}
	}()
}

var NimSlashCommand = discord.SlashCommandCreate{
	Name:        "nim",
	Description: "nim game where you take rocks until you get 0",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "choose your enemy",
			Required:    true, //make it false later when the bot is added

		},
	},
}

/*

Later add ai so user can battle with the ai


*/
