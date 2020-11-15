package gossip

import (
	"github.com/raidancampbell/gossip/data"
	"gopkg.in/sorcix/irc.v2"
)

var pingPong = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PING
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.PONG,
			Params:  msg.Params,
		}
		return false
	},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "pingPong",
	},
}

