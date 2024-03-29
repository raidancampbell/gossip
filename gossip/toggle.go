package gossip

import (
	"fmt"
	"github.com/raidancampbell/gossip/data"
	"gopkg.in/sorcix/irc.v2"
	"strings"
)


var toggle = &ComposedTrigger{
	subTriggers: []Trigger{triggerToggle, triggerStatus},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "toggle",
	},
}

// toggle triggers on !toggle, !enable, !disable
var triggerToggle = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && (strings.HasPrefix(msg.Params[1], "!toggle") || strings.HasPrefix(msg.Params[1], "!enable") || strings.HasPrefix(msg.Params[1], "!disable"))  && msg.Name == g.cfg.OwnerNick
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		words := strings.Split(msg.Params[1], " ")
		if len(words) != 2 {
			return true
		}

		var t []string
		for i := range g.triggers {
			t = append(t, g.triggers[i].GetMeta().Name)
			if strings.EqualFold(g.triggers[i].GetMeta().Name, words[1]) {
				var nextDisabled bool
				if strings.EqualFold(words[0], "!toggle") {
					nextDisabled = !g.triggers[i].GetMeta().Disabled
				} else if strings.EqualFold(words[0], "!enable") {
					nextDisabled = false
				} else if strings.EqualFold(words[0], "!disable") {
					nextDisabled = true
				}

				g.msgChan <- &irc.Message{
					Command: irc.PRIVMSG,
					Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("%s: was %s, now %s", g.triggers[i].GetMeta().Name, disabledToStr(g.triggers[i].GetMeta().Disabled), disabledToStr(nextDisabled))},
				}
				g.triggers[i].GetMeta().Disabled = nextDisabled

				g.db.Model(&data.TriggerMeta{}).Where(&data.TriggerMeta{Name: g.triggers[i].GetMeta().Name}).Update("disabled", g.triggers[i].GetMeta().Disabled)
				return false
			}
		}

		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("unknown trigger '%s', known triggers: [%s]", words[1], strings.Join(t, ", "))},
		}
		return false
	},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "triggerToggle",
	},
}

func disabledToStr(b bool) string {
	if b {
		return "disabled"
	}
	return "enabled"
}

var triggerStatus = &SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!status"
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		var enabled []string
		var disabled []string

		for i := range g.triggers {
			if g.triggers[i].GetMeta().Disabled {
				disabled = append(disabled, g.triggers[i].GetMeta().Name)
			} else {
				enabled = append(enabled, g.triggers[i].GetMeta().Name)
			}
		}

		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("Enabled: [%s]", strings.Join(enabled, ", "))},
		}
		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("Disabled: [%s]", strings.Join(disabled, ", "))},
		}
		return false
	},
	meta: &data.TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "triggerToggle",
	},
}