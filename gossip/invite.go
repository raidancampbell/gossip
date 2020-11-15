package gossip

import (
	"github.com/raidancampbell/gossip/data"
	"gopkg.in/sorcix/irc.v2"
)

// on /invite, join the desired channel
var invite = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.INVITE
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.JOIN,
			Params:  []string{msg.Params[1]},
		}
		return false
	},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "invite",
	},
}
