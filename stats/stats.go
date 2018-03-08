package stats

import (
	"time"
)

type Stats struct {
	// General bot stats.
	guilds      map[string]*Guild
	connectedAt time.Time

	// Actions queue.
	actionsQueue chan Action

	// User defined stats.
	bag *Bag
}

func NewStats(actionsQueueSize int) *Stats {
	return &Stats{
		guilds:       make(map[string]*Guild),
		connectedAt:  time.Now(),
		actionsQueue: make(chan Action, actionsQueueSize),
		bag:          NewBag(),
	}
}

func (s *Stats) Guilds() map[string]*Guild {
	return s.guilds
}

func (s *Stats) ConnectedAt() time.Time {
	return s.connectedAt
}

func (s *Stats) Bag() *Bag {
	return s.bag
}

func (s *Stats) PushAction(a Action) {
	s.actionsQueue <- a
}

func (s *Stats) ProcessActions() {
	for a := range s.actionsQueue {
		switch a.Name() {
		case "create_guild":
			s.doCreateGuild(a)
		case "delete_guild":
			s.doDeleteGuild(a)
		}
	}
}

func (s *Stats) doCreateGuild(a Action) {
	ca, converted := a.(*CreateGuildAction)
	if !converted {
		// error!
		return
	}

	_, found := s.guilds[ca.GuildId]
	if !found {
		s.guilds[ca.GuildId] = &Guild{
			id:     ca.GuildId,
			name:   ca.GuildName,
			region: ca.GuildRegion,
		}
	}
}

func (s *Stats) doDeleteGuild(a Action) {
	ca, converted := a.(*DeleteGuildAction)
	if !converted {
		// error!
		return
	}

	_, found := s.guilds[ca.GuildId]
	if found {
		delete(s.guilds, ca.GuildId)
	}
}
