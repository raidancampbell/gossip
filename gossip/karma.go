package gossip

import (
	"fmt"
	"github.com/raidancampbell/gossip/data"
	"gopkg.in/sorcix/irc.v2"
	"strings"
)


var karmaCounter = SyncTrigger{
	Cond: func(bot *Bot, msg *irc.Message) (shouldApply bool) {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && (strings.Contains(msg.Params[1], "++") || strings.Contains(msg.Params[1], "--"))
	},
	Act: func(g *Bot, msg *irc.Message) (shouldContinue bool) {
		var words []string
		incrBy := -1
		for _, word := range strings.Split(msg.Params[1], " ") {
			// if we're starting a (karma this)++, clobber any existing state
			// I am NOT about to handle (multiple (nests of)++ karma stuff)--
			if strings.HasPrefix(word, "(") || strings.HasPrefix(word, "\"") {
				words = []string{strings.TrimLeft(word, "(\"")}
			} else if len(words) > 0 && !(strings.HasSuffix(word, "++") || strings.HasSuffix(word, "--")) {
				words = append(words, word)
			} else if strings.HasSuffix(word, "++") || strings.HasSuffix(word, "--") {
				if strings.HasSuffix(word, "++") {
					incrBy = 1
				}

				trimmedWord := strings.TrimRight(word, "+-")
				if strings.HasSuffix(trimmedWord, ")") || strings.HasSuffix(trimmedWord, "\"") {
					trimmedWord = strings.TrimRight(trimmedWord, ")\"")
					words = append(words, trimmedWord)
				} else {
					words = []string{trimmedWord}
				}
				break
			}
		}
		if len(words) == 0 {
			// the match is pretty rough
			// e.g. foo++bar will match.
			//TODO: matching should be a regex, action is too complicated
			return
		}
		token := strings.Join(words, " ")
		k := &data.Karma{
			Object:   token,
			Value:    0,
			Location: mirrorMsg(g, msg),
		}

		g.db.Model(k).Where(k).First(k)
		k.Value += incrBy
		g.db.Model(k).Save(k)

		return true
	},
}

var KarmaBest = SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!karma best"
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		var karma []data.Karma

		// in the karma table, sorted by value descending, get the top 5, where we match the channel
		g.db.Model(&data.Karma{}).Order("value desc").Limit(5).Where(&data.Karma{
			Location: mirrorMsg(g, msg),
		}).Find(&karma)

		for _, k := range karma{
			g.msgChan <- &irc.Message{
				Command: irc.PRIVMSG,
				Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("%d: %s", k.Value, k.Object)},
			}
		}
		return false
	},
}

var KarmaWorst = SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && msg.Params[1] == "!karma worst"
	},
	Act: func(g *Bot, msg *irc.Message) bool {
		var karma []data.Karma

		// in the karma table, sorted by value descending, get the top 5, where we match the channel
		g.db.Model(&data.Karma{}).Order("value asc").Limit(5).Where(&data.Karma{
			Location: mirrorMsg(g, msg),
		}).Find(&karma)

		for _, k := range karma{
			g.msgChan <- &irc.Message{
				Command: irc.PRIVMSG,
				Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("%d: %s", k.Value, k.Object)},
			}
		}
		return false
	},
}