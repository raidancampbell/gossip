package gossip

import (
	"fmt"
	"github.com/gregdel/pushover"
	"github.com/raidancampbell/gossip/conf"
	"github.com/raidancampbell/gossip/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
	"strings"
)

// sends push notifications on highlighted keywords
type pushTrigger struct {
	a           *pushover.Pushover
	r           *pushover.Recipient
	highlightOn []string
	meta        *data.TriggerMeta
}

func NewPush(cfg *conf.Cfg) Trigger {
	app := pushover.New(cfg.Triggers.Push.APIKey)
	recip := pushover.NewRecipient(cfg.Triggers.Push.RecipientKey)
	cfg.Triggers.Push.HighlightOn = append(cfg.Triggers.Push.HighlightOn, cfg.OwnerNick)
	return &pushTrigger{
		a:           app,
		r:           recip,
		highlightOn: cfg.Triggers.Push.HighlightOn,
		meta: &data.TriggerMeta{
			Disabled: false,
			Priority: 0,
			Name:     "push",
		},
	}
}

func (p pushTrigger) Condition(_ *Bot, msg *irc.Message) (shouldApply bool) {
	return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && ContainsAny(msg.Params[1], p.highlightOn...)
}

func (p pushTrigger) Action(_ *Bot, msg *irc.Message) (shouldContinue bool) {
	notif := &pushover.Message{
		Message: fmt.Sprintf("New IRC Highlight in %s\n%s: %s", msg.Params[0], msg.Name, msg.Params[1]),
		Title:   "IRC Highlight",
	}
	resp, err := p.a.SendMessage(notif, p.r)
	if err != nil {
		logrus.WithError(err).WithField("pushover-response", resp).Error("Push notification to pushover failed")
	}
	logrus.WithField("pushover-response", resp).Info("Push notification sent successfully")
	return true
}

func (p pushTrigger) GetMeta() *data.TriggerMeta {
	return p.meta
}

func (p pushTrigger) Meta(meta *data.TriggerMeta) {
	p.meta = meta
}

func ContainsAny(haystack string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}
