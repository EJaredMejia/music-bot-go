package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/music-formatter/queue"
)

const (
	BR = "br"
)

type CommandParams struct {
	Queues         queue.DynamicQueues
	TextMessage    string
	Discord        *discordgo.Session
	DiscordMessage *discordgo.MessageCreate
	Flags          []string
}
