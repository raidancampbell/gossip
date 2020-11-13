package gossip

import (
	"gopkg.in/sorcix/irc.v2"
	"strings"
)

// 7 bits or bust
// also don't ever say the word 'moist'
var hiss = SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		if msg.Command != irc.PRIVMSG {
			return false
		}
		for _, c := range []rune(msg.Params[1]) {
			if c > 127 {
				return true
			}
		}
		if strings.Contains(strings.ToLower(msg.Params[1]), "moist") {
			return true
		}
		return false
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{mirrorMsg(g, msg), "hisss"},
		}
		return false
	},
}