package gossip

import (
	"github.com/raidancampbell/gossip/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
	"os"
	"time"
)

// exit on !die
var die = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!die" && msg.Name == g.cfg.OwnerNick
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		logrus.Info("Exiting...")
		g.msgChan <- &irc.Message{
			Command: irc.QUIT,
			Params:  []string{"goodbye"},
		}
		time.Sleep(100 * time.Millisecond)
		g.c.Close()
		os.Exit(0)
		return false
	},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "die",
	},
}