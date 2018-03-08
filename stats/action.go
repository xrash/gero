package stats

import (
	"time"
)

type Action interface {
	Name() string
}

type CreateGuildAction struct {
	GuildId     string
	GuildName   string
	GuildRegion string
	CreatedAt   time.Time
}

func (a *CreateGuildAction) Name() string {
	return "create_guild"
}

type DeleteGuildAction struct {
	GuildId string
}

func (a *DeleteGuildAction) Name() string {
	return "delete_guild"
}
