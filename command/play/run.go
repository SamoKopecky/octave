package play

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/octave/command"
	"github.com/lukasl-dev/waterlink/v2/track/query"
)

func run(deps Deps) command.Handler {
	return func(s *discordgo.Session, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
		q := command.Option(evt.ApplicationCommandData().Options, "query")

		res, err := deps.Client.LoadTracks(query.Of(q.StringValue()))
		switch {
		case err != nil:
			return command.ErrorResponse(errors.New("failed to load track"))
		case len(res.Tracks) == 0:
			return command.ErrorResponse(errors.New("no tracks found"))
		}
		tr := res.Tracks[0]

		channelID := command.MemberChannel(s, evt.GuildID, evt.Member.User.ID)
		if channelID == "" {
			return command.ErrorResponse(errors.New("user is not a member of any voice channel"))
		}

		err = s.ChannelVoiceJoinManual(evt.GuildID, channelID, false, true)
		if err != nil {
			return command.ErrorResponse(errors.New("failed to join voice channel"))
		}

		err = deps.Conn.Guild(snowflake.MustParse(evt.GuildID)).PlayTrack(tr)
		if err != nil {
			return command.ErrorResponse(errors.New("failed to play track"))
		}

		s.ChannelMessageSend(evt.ChannelID, fmt.Sprintf("Playing **%s** with link %s", tr.Info.Title, tr.Info.URI))
		log.Printf("Playing: %+v\n", tr.Info)

		return &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "▶️",
							},
							Label: fmt.Sprintf("Now playing %s", tr.Info.Title),
							Style: discordgo.LinkButton,
							URL:   tr.Info.URI,
						},
					},
				},
			},
		}
	}
}
