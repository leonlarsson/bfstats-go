package commandhandlers

import (
	"errors"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/leonlarsson/bfstats-bot-go/canvas"
	"github.com/leonlarsson/bfstats-bot-go/canvasdatashapes"
	create "github.com/leonlarsson/bfstats-bot-go/create/bf2042"
	"github.com/leonlarsson/bfstats-bot-go/datafetchers/bf2042datafetcher"
	"github.com/leonlarsson/bfstats-bot-go/localization"
	"github.com/leonlarsson/bfstats-bot-go/shared"
	"github.com/leonlarsson/bfstats-bot-go/utils"
)

// HandleBF2042VehiclesCommand handles the bf2042 vehicles command.
func HandleBF2042VehiclesCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate, loc localization.LanguageLocalizer) error {
	username, platform, usernameFailedValidation := utils.GetStatsCommandArgs(session, interaction, &loc)
	if usernameFailedValidation {
		return errors.New("username failed validation")
	}

	data, err := bf2042datafetcher.FetchBF2042VehiclesData(platform, username)
	if err != nil {
		return err
	}

	if len(data.Data) < 9 {
		return errors.New(loc.Translate("messages/not_enough_vehicles", map[string]string{"vehicles": string(rune(len(data.Data)))}))
	}

	// Sort the vehicles by kills
	sort.Slice(data.Data, func(i, j int) bool {
		return data.Data[i].Stats.Kills.Value > data.Data[j].Stats.Kills.Value
	})

	// Build the vehicles slice
	var vehicles []canvasdatashapes.Stat
	for _, vehicle := range data.Data {
		vehicleSlice := canvasdatashapes.Stat{
			Name:  strings.TrimSpace(vehicle.Metadata.Name),
			Value: loc.Translate("stats/title/x_kills_short", map[string]string{"kills": loc.FormatInt(vehicle.Stats.Kills.Value)}),
			Extra: loc.Translate("stats/title/x_kpm_and_time", map[string]string{"kpm": loc.FormatFloat(vehicle.Stats.KillsPerMinute.Value), "time": vehicle.Stats.TimePlayed.DisplayValue}),
		}
		vehicles = append(vehicles, vehicleSlice)
	}

	// Create the image
	imageData := canvasdatashapes.BF2042VehiclesCanvasData{
		BaseData: canvasdatashapes.BaseData{
			Identifier: "BF2042-001",
			Username:   username,
			Platform:   utils.TRNPlatformNameToInt(platform),
			Meta: canvasdatashapes.Meta{
				Game:    "Battlefield 2042",
				Segment: "Vehicles",
			},
		},
		Vehicles: vehicles,
	}

	c, _ := create.CreateBF2042VehiclesImage(imageData, shared.SolidBackground)
	imgBuf := canvas.CanvasToBuffer(c)

	// Edit the response
	session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Files: []*discordgo.File{
			{
				Name:        "vehicles.png",
				ContentType: "image/png",
				Reader:      imgBuf,
			},
		},
	})

	return nil
}
