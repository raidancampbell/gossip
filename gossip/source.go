package gossip

import "gopkg.in/sorcix/irc.v2"

// on !source link the source code
var source = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!source"
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{mirrorMsg(g, msg), "https://github.com/raidancampbell/gossip"},
		}
		return false
	},
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "source",
	},
}
