package gossip

import "gopkg.in/sorcix/irc.v2"

// pattern stolen from https://github.com/whyrusleeping/hellabot/blob/master/hellabot.go
type StatelessTrigger struct {
	// Returns true if this trigger applies to the passed in message
	Cond func(*Bot, *irc.Message) (shouldApply bool)

	// The action to perform if Cond is true
	// return true if processing should continue
	Act func(*Bot, *irc.Message) (shouldContinue bool)
}

func (t StatelessTrigger) Condition(b *Bot, msg *irc.Message) (shouldApply bool) {
	return t.Cond(b, msg)
}
func (t StatelessTrigger) Action(b *Bot, msg *irc.Message) (shouldContinue bool) {
	return t.Act(b, msg)
}

type Trigger interface {
	Condition(*Bot, *irc.Message) (shouldApply bool)
	Action(*Bot, *irc.Message) (shouldContinue bool)
}

var pingPong = StatelessTrigger{
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
}

// join desired channels on startup
var joinChans = StatelessTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.RPL_WELCOME
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		g.joinChannels.Do(func() {
			for _, chn := range g.channels {
				g.msgChan <- &irc.Message{
					Command: irc.JOIN,
					Params:  []string{chn},
				}
			}
		})
		return true
	},
}

// on /invite, join the desired channel
var invite = StatelessTrigger{
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
}

// on !ping, give pong!
var userPingPong = StatelessTrigger{
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
}