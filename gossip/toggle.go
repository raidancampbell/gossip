package gossip

import (
	"fmt"
	"gopkg.in/sorcix/irc.v2"
	"strings"
)

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
				var nextState bool
				if strings.EqualFold(words[0], "!toggle") {
					nextState = !g.triggers[i].GetMeta().Disabled
				} else if strings.EqualFold(words[0], "!enable") {
					nextState = true
				} else if strings.EqualFold(words[0], "!false") {
					nextState = false
				}
				g.msgChan <- &irc.Message{
					Command: irc.PRIVMSG,
					Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("%s: was %s, now %s", g.triggers[i].GetMeta().Name, disabledToStr(g.triggers[i].GetMeta().Disabled), disabledToStr(nextState))},
				}
				g.triggers[i].GetMeta().Disabled = nextState
				return false
			}
		}

		g.msgChan <- &irc.Message{
			Command: irc.PRIVMSG,
			Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("unknown trigger '%s', known triggers: [%s]", words[1], strings.Join(t, ", "))},
		}
		return false
	},
	meta: TriggerMeta{
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
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "triggerToggle",
	},
}