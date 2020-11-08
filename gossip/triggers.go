package gossip

import "gopkg.in/sorcix/irc.v2"

// pattern stolen from https://github.com/whyrusleeping/hellabot/blob/master/hellabot.go
type Trigger struct {
	// Returns true if this trigger applies to the passed in message
	Condition func(*Bot, *irc.Message) (shouldApply bool)

	// The action to perform if Condition is true
	// return true if processing should continue
	Action func(*Bot, *irc.Message) (shouldContinue bool)
}

var pingPong = Trigger{
	Condition: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PING
	},
	Action: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.PONG,
			Params:  msg.Params,
		}
		return false
	},
}

// join desired channels on startup
var joinChans = Trigger{
	Condition: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.RPL_WELCOME
	},
	Action: func(g *Bot, msg *irc.Message) bool {
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
var invite = Trigger {
	Condition: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.INVITE
	},
	Action: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.JOIN,
			Params:  []string{msg.Params[1]},
		}
		return false
	},
}

// on !ping, give pong!
var userPingPong = Trigger {
	Condition: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!ping"
	},
	Action: func(g *Bot, msg *irc.Message) bool {
		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{msg.Params[0], "pong!"},
		}
		return false
	},
}