package gero

import (
	"time"
)

type StatusManager struct {
	started                 bool
	status                  string
	bot                     *Bot
	exitUpdateStatusChannel chan bool
}

func NewStatusManager() *StatusManager {
	return &StatusManager{
		exitUpdateStatusChannel: make(chan bool, 1),
	}
}

func (sm *StatusManager) WithBot(bot *Bot) *StatusManager {
	if !sm.started {
		sm.started = true
		go sm.processStatusChanges()
	}

	sm.bot = bot

	return sm
}

func (sm *StatusManager) Set(s string) {
	sm.status = s
}

func (sm *StatusManager) Exit() {
	sm.exitUpdateStatusChannel <- true
}

func (sm *StatusManager) processStatusChanges() {
	for {
		select {
		case <-sm.exitUpdateStatusChannel:
			return
		default:
			if err := sm.bot.session.UpdateStatus(0, sm.status); err != nil {
				sm.bot.Logger().Error("Error updating status to %s: %v", sm.status, err)
			}
		}
		time.Sleep(time.Second * 60)
	}
}
