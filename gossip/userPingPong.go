package gossip

import (
	"github.com/raidancampbell/gossip/data"
	"gopkg.in/sorcix/irc.v2"
)

// on !ping, give pong!
var userPingPong = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!ping"
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{mirrorMsg(g, msg), "pong!"},
		}
		return false
	},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "userPingPong",
	},
}
