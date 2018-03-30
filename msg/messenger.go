package msg

import (
	"github.com/bwmarrin/discordgo"
)

type Messenger struct {
	session *discordgo.Session
	logger  Logger
	color   int
}

func NewMessenger(session *discordgo.Session, logger Logger, color int) *Messenger {
	return &Messenger{
		session: session,
		logger:  logger,
		color:   color,
	}
}

func (m *Messenger) Send(channelId, content string) {
	go func() {
		_, err := m.session.ChannelMessageSend(channelId, content)
		if err != nil {
			m.logger.Error("Error sending message: %v", err)
		}
	}()
}

func (m *Messenger) SendTemplate(channelId, template string, data interface{}) {
	go func() {
		content, err := ProcessTemplate(template, data)
		if err != nil {
			m.logger.Error("Error processing template: %v", err)
		}

		m.Send(channelId, content)
	}()
}

func (m *Messenger) SendEmbed(channelId, content string) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       m.color,
		Description: content,
	}

	go func() {
		_, err := m.session.ChannelMessageSendEmbed(channelId, embed)
		if err != nil {
			m.logger.Error("Error sending embed message: %v", err)
		}
	}()
}

func (m *Messenger) SendEmbedTemplate(channelId, template string, data interface{}) {
	go func() {
		content, err := ProcessTemplate(template, data)
		if err != nil {
			m.logger.Error("Error processing template: %v", err)
		}

		m.SendEmbed(channelId, content)
	}()
}
