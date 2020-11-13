package gossip

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
	"strings"
)

// rename on !rename
var rename = SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!rename" && len(strings.Split(msg.Params[1], " ")) == 2
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		logrus.Infof("renaming to %s", strings.Split(msg.Params[1], " ")[1])
		g.msgChan <- &irc.Message{
			Command: irc.NICK,
			Params:  []string{strings.Split(msg.Params[1], " ")[1]},
		}
		return false
	},
}