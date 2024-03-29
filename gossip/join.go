package gossip

import (
	"github.com/raidancampbell/gossip/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
	"sync"
	"time"
)

// join desired channels on startup
type joinChannels struct {
	channels []string
	o *sync.Once
	meta *data.TriggerMeta
}

func NewJoin(channels []string) Trigger {
	return &joinChannels{
		channels: channels,
		o:        &sync.Once{},
		meta:     &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "joinChans",
	},
	}
}

func (j joinChannels) Condition(_ *Bot, msg *irc.Message) (shouldApply bool) {
	return msg.Command == irc.RPL_WELCOME
}

func (j joinChannels) Action(g *Bot, _ *irc.Message) (shouldContinue bool) {
	j.o.Do(func() {
		go func() {
			if len(g.cfg.Network.OnConnect) > 0 {
				waitTime := 10 * time.Second
				logrus.Infof("waiting %s before joining channels to allow onConnect command to run...", waitTime.String())
				time.Sleep(waitTime)
			}
			for _, chn := range j.channels {
				g.msgChan <- &irc.Message{
					Command: irc.JOIN,
					Params:  []string{chn},
				}
			}
		}()
	})
	return true
}

func (j joinChannels) GetMeta() *data.TriggerMeta {
	return j.meta
}

func (j joinChannels) Meta(meta *data.TriggerMeta) {
	j.meta = meta
}