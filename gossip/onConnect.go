package gossip

import (
	"github.com/raidancampbell/gossip/data"
	"gopkg.in/sorcix/irc.v2"
	"time"
)

var onConnect = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.RPL_WELCOME
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		time.Sleep(1 * time.Second)
		for _, rawCMD := range g.cfg.Network.OnConnect {
			g.msgChan <- irc.ParseMessage(rawCMD)
		}
		return true
	},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "onConnect",
	},
}
