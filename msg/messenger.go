package msg

import (
	"github.com/bwmarrin/discordgo"
)

type Messenger struct {
	session          *discordgo.Session
	logger           Logger
	color            int
	sentMessages     *MessengerBag
	receivedMessages *MessengerBag
}

func NewMessenger(session *discordgo.Session, logger Logger, color int) *Messenger {
	return &Messenger{
		session:          session,
		logger:           logger,
		color:            color,
		sentMessages:     NewMessengerBag(48, 16),
		receivedMessages: NewMessengerBag(48, 16),
	}
}

func (ms *Messenger) BuildEmbed(content string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       ms.color,
		Description: content,
	}
}

func (ms *Messenger) Send(channelId, content string) {
	go func() {
		m, err := ms.session.ChannelMessageSend(channelId, content)
		if err != nil {
			ms.logger.Error("Error sending message: %v", err)
		} else {
			ms.sentMessages.Capture(channelId, m.ID)
		}
	}()
}

func (ms *Messenger) SendTemplate(channelId, template string, data interface{}) {
	go func() {
		content, err := ProcessTemplate(template, data)
		if err != nil {
			ms.logger.Error("Error processing template: %v", err)
		}

		ms.Send(channelId, content)
	}()
}

func (ms *Messenger) SendEmbed(channelId, content string) {
	embed := ms.BuildEmbed(content)

	go func() {
		m, err := ms.session.ChannelMessageSendEmbed(channelId, embed)
		if err != nil {
			ms.logger.Error("Error sending embed message: %v", err)
		} else {
			ms.sentMessages.Capture(channelId, m.ID)
		}
	}()
}

func (ms *Messenger) SendEmbedTemplate(channelId, template string, data interface{}) {
	go func() {
		content, err := ProcessTemplate(template, data)
		if err != nil {
			ms.logger.Error("Error processing template: %v", err)
		}

		ms.SendEmbed(channelId, content)
	}()
}

func (ms *Messenger) ReceivedMessages() *MessengerBag {
	return ms.receivedMessages
}

func (ms *Messenger) SentMessages() *MessengerBag {
	return ms.sentMessages
}

func (ms *Messenger) Clean(channelId string) (int, error) {
	var total int
	var err error
	var s, r map[string][]string

	ms.sentMessages.WithBag(func(sent map[string][]string) map[string][]string {
		ms.receivedMessages.WithBag(func(recv map[string][]string) map[string][]string {
			s, r, total, err = ms.doClean(channelId, sent, recv)
			return r
		})
		return s
	})

	return total, err
}

func (ms *Messenger) doClean(channelId string, sent, recv map[string][]string) (map[string][]string, map[string][]string, int, error) {
	sqtt := len(sent[channelId])
	rqtt := len(recv[channelId])
	total := sqtt + rqtt

	if total < 1 {
		return sent, recv, 0, nil
	}

	messages := make([]string, 0)
	messages = append(messages, sent[channelId]...)
	messages = append(messages, recv[channelId]...)

	err := ms.session.ChannelMessagesBulkDelete(channelId, messages)

	sent[channelId] = make([]string, 0)
	recv[channelId] = make([]string, 0)

	return sent, recv, total, err
}
