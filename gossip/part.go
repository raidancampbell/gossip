package gossip

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
)

// leave on !part
var part = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!part" && msg.Name == g.cfg.OwnerNick
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		logrus.Info("leaving...")
		g.msgChan <- &irc.Message{
			Command: irc.PART,
			Params:  []string{msg.Params[0]},
		}
		return false
	},
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "part",
	},
}