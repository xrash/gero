package msg

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type Messenger struct {
	session *discordgo.Session
	logger  Logger
	color   int

	sentLock *sync.Mutex
	sent     map[string][]string
}

func NewMessenger(session *discordgo.Session, logger Logger, color int) *Messenger {
	return &Messenger{
		session:  session,
		logger:   logger,
		color:    color,
		sentLock: &sync.Mutex{},
		sent:     make(map[string][]string),
	}
}

func (ms *Messenger) Send(channelId, content string) {
	go func() {
		m, err := ms.session.ChannelMessageSend(channelId, content)
		if err != nil {
			ms.logger.Error("Error sending message: %v", err)
		} else {
			ms.capture(channelId, m.ID)
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
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       ms.color,
		Description: content,
	}

	go func() {
		m, err := ms.session.ChannelMessageSendEmbed(channelId, embed)
		if err != nil {
			ms.logger.Error("Error sending embed message: %v", err)
		} else {
			ms.capture(channelId, m.ID)
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

func (ms *Messenger) Clean(channelId string) (int, error) {
	ms.sentLock.Lock()
	defer ms.sentLock.Unlock()

	qtt := len(ms.sent[channelId])

	if qtt < 1 {
		return -1, nil
	}

	err := ms.session.ChannelMessagesBulkDelete(channelId, ms.sent[channelId])
	if err != nil {
		ms.logger.Error("Error processing template: %v", err)
	}

	ms.sent[channelId] = make([]string, 0)

	return qtt, err
}

func (ms *Messenger) capture(cid, mid string) {
	ms.sentLock.Lock()
	defer ms.sentLock.Unlock()

	maxListSize := 80
	minListSize := maxListSize / 2

	list := ms.sent[cid]
	if list == nil {
		list = make([]string, 0)
	}

	list = append(list, mid)

	if len(list) > maxListSize {
		list = list[len(list)-minListSize:]
	}

	ms.sent[cid] = list
}
