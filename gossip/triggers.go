package gossip

import "gopkg.in/sorcix/irc.v2"

// pattern stolen from https://github.com/whyrusleeping/hellabot/blob/master/hellabot.go
type SyncTrigger struct {
	// Returns true if this trigger applies to the passed in message
	Cond func(*Bot, *irc.Message) (shouldApply bool)

	// The action to perform if Cond is true
	// return true if processing should continue
	Act func(*Bot, *irc.Message) (shouldContinue bool)

	meta TriggerMeta
}

// metadata for triggers, extracted for ease of iteration
type TriggerMeta struct {
	Disabled bool
	Priority int // maybe a comparator? maybe don't need it at all?
	Name string
}

func (t SyncTrigger) Condition(b *Bot, msg *irc.Message) (shouldApply bool) {
	return t.Cond(b, msg)
}
func (t SyncTrigger) Action(b *Bot, msg *irc.Message) (shouldContinue bool) {
	return t.Act(b, msg)
}
func (t SyncTrigger) GetMeta() TriggerMeta {
	return t.meta
}

type Trigger interface {
	Condition(*Bot, *irc.Message) (shouldApply bool)
	Action(*Bot, *irc.Message) (shouldContinue bool)
	GetMeta() TriggerMeta
}

var pingPong = SyncTrigger{
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
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "pingPong",
	},
}

// join desired channels on startup
var joinChans = SyncTrigger{
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
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "joinChans",
	},

}

// on /invite, join the desired channel
var invite = SyncTrigger{
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
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "invite",
	},
}

// on !ping, give pong!
var userPingPong = SyncTrigger{
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
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "userPingPong",
	},
}