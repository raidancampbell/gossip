package gossip

import "gopkg.in/sorcix/irc.v2"

func mirrorMsg(g *Bot, msg *irc.Message) string {
	if len(msg.Params) == 0 {
		//oh no
		return ""
	}
	if msg.Params[0] == g.nick {
		return msg.Name
	}
	return msg.Params[0]
}