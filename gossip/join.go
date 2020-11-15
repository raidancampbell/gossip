package gossip

import (
	"github.com/raidancampbell/gossip/data"
	"gopkg.in/sorcix/irc.v2"
)

// join desired channels on startup
//TODO: refactor this into a different impl, and initialize it with the channel array
//this removes the need for keeping the array in the bot struct
var joinChans = &SyncTrigger{
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
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "joinChans",
	},
}