package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/leonlarsson/bfstats-go/internal/bot/commands"
	"github.com/leonlarsson/bfstats-go/internal/datafetchers/trndatafetcher"
)

func HandleInteractionCreateEvent(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Autocomplete
	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		cmdData := interaction.ApplicationCommandData()
		options := commands.ParseOptions(cmdData.Options)
		game := cmdData.Name
		username := options["username"].StringValue()
		platform := options["platform"].StringValue()

		var responseData []*discordgo.ApplicationCommandOptionChoice

		// If fewer than 3 characters. Else fetch data from TRN.
		if len(username) < 3 {
			// If no input, respond with user's display name. Else, respond with the input.
			if len(username) == 0 {
				responseData = []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  interaction.Member.DisplayName(),
						Value: interaction.Member.DisplayName(),
					},
				}
			} else {
				responseData = []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  username,
						Value: username,
					},
				}
			}
		} else {
			// Fetch data
			data, err := trndatafetcher.FetchTRNSearchData(game, platform, username)
			if err != nil {
				return
			}

			// Build the response data for the first 25 users
			for i, user := range data.Data {
				if i >= 25 {
					break
				}
				responseData = append(responseData, &discordgo.ApplicationCommandOptionChoice{
					Name:  user.PlatformUserIdentifier,
					Value: user.PlatformUserIdentifier,
				})
			}
		}

		// Respond with the data
		s.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: responseData,
			},
		})
	}

	// Chat input command
	// TODO: Just prototyping. Add actual command handling based on the command name and subcommand if needed.
	if interaction.Type == discordgo.InteractionApplicationCommand {
		cmdData := interaction.ApplicationCommandData()
		options := commands.ParseOptions(cmdData.Options)
		println("Command: ", cmdData.Name)
		println("Subcommand: ", commands.GetOptionStringValue(options, "subcommand"))
		println("Username: ", commands.GetOptionStringValue(options, "username"))
		println("Platform: ", commands.GetOptionStringValue(options, "platform"))
		println("Format: ", commands.GetOptionStringValue(options, "format"))
		println("Language: ", commands.GetOptionStringValue(options, "language"))
		println("PoemGPT: ", commands.GetOptionBoolValue(options, "poem_gpt"))
	}
}
