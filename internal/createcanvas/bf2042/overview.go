package bf2042

import (
	"github.com/leonlarsson/bfstats-go/internal/canvas"
	"github.com/leonlarsson/bfstats-go/internal/canvas/shapes"
	"github.com/leonlarsson/bfstats-go/internal/shared"
	"github.com/leonlarsson/bfstats-go/internal/utils"
	core "github.com/tdewolff/canvas"
)

func CreateBF2042OverviewImage(data shapes.GenericRegularData, style shared.BackgroundFormat) (*core.Canvas, *core.Context) {
	c, ctx := canvas.BuildBaseCanvas("BF2042", data.BaseData, shared.RegularSkeletonType)

	// canvas.DrawTimePlayed(ctx, data.Stats.TimePlayed)

	canvas.DrawStat1(ctx, data.Stats.L1)
	canvas.DrawStat2(ctx, data.Stats.L2)

	canvas.DrawStat3(ctx, data.Stats.L3)
	canvas.DrawStat4(ctx, data.Stats.L4)

	canvas.DrawBestClassImage(ctx, "BF2042", data.Stats.L5.Value)

	// 1. If the best class is a base class, draw the text with room for the image
	// 2. If the best class is "Unknown", draw the multi-kills stat
	// 3. If the best class is anything else ("BF3 Engineer, etc."), draw a regular stat
	if utils.IsBaseBF2042Class("BF2042", data.Stats.L5.Value) {
		canvas.DrawStat5BestClass(ctx, data.Stats.L5)
	} else if data.Stats.L5.Value == "Unknown" {
		canvas.DrawStat5(ctx, data.Stats.L5Fallback)
	} else {
		canvas.DrawStat5(ctx, data.Stats.L5)
	}

	canvas.DrawStat6(ctx, data.Stats.L6)

	canvas.DrawRightStat1(ctx, data.Stats.R1)
	canvas.DrawRightStat2(ctx, data.Stats.R2)
	canvas.DrawRightStat3(ctx, data.Stats.R3)
	canvas.DrawRightStat4Rank(ctx, data.Stats.R4)

	return c, ctx
}
