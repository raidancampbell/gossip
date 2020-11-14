package gossip

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"gopkg.in/sorcix/irc.v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	maxHTMLResponseBytes = 1024 * 1024 * 5 // 5 MB
)

// prints the html title text of any URLs within a message
var htmlTitle = SyncTrigger{
	Cond: func(g *Bot, msg *irc.Message) bool {
		if msg.Command != irc.PRIVMSG {
			return false
		}
		w := extractURLWord(msg.Params[1])
		return w != ""
	},

	Act: func(g *Bot, msg *irc.Message) (shouldContinue bool) {
		shouldContinue = true // always
		var u *url.URL
		for _, word := range strings.Split(msg.Params[1], " ") {
			if strings.Contains(word, "://"){
				tmpURL, err := url.Parse(strings.Trim(word, ":,!.<>"))
				if err == nil {
					u = tmpURL
					break
				}
			}
		}
		// not possible
		if u == nil {
			return
		}
		logrus.Debugf("retrieving title for URL '%s'", u.String())

		r, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			logrus.WithError(err).Error("failed to create HTTP request to endpoint %+v", u)
			return
		}

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			logrus.WithError(err).Error("failed to complete HTTP request to endpoint %+v", u)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxHTMLResponseBytes))
		if err != nil {
			logrus.WithError(err).Error("failed to read %d bytes from %+v", maxHTMLResponseBytes, u)
			return
		}
		logrus.Debugf("URL %s had an HTML body of length %d", u.String(), len(body))
		tok := html.NewTokenizer(bytes.NewReader(body))
		for {
			tokType := tok.Next()
			if tokType == html.ErrorToken {
				if tok.Err() != io.EOF {
					logrus.WithError(tok.Err()).Error("error while parsing HTML response")
				}
				break
			}

			t := tok.Token()
			if t.Data == "title" {
				nextType := tok.Next()
				if nextType == html.TextToken {
					g.msgChan <- &irc.Message{
						Command: irc.PRIVMSG,
						Params:  []string{mirrorMsg(g, msg), tok.Token().Data},
					}
					break
				}
			}
		}
		return
	},
	meta: TriggerMeta{
		Disabled: false,
		Priority: 0,
		Name:     "HTMLTitle",
	},
}

func extractURLWord(s string) string {
	for _, word := range strings.Split(s, " ") {
		if strings.Contains(word, "://"){
			_, err := url.Parse(strings.Trim(word, ":,!.<>"))
			if err == nil {
				return word
			}
		}
	}
	return ""
}