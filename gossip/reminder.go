package gossip

import (
	"fmt"
	"github.com/raidancampbell/gossip/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reminderPattern = regexp.MustCompile(`(?i)^!remind (?P<Target>(me|[a-zA-Z]+) )?in ((?P<Seconds>\d+ ?(seconds?|secs?))|(?P<Minutes>\d+ ?(minutes?|mins?))|(?P<Hours>\d+ ?(hours?|hrs?))|(?P<Days>\d+ ?days?)|(?P<Weeks>\d+ ?weeks?)|(?P<Months>\d+ ?(months?|mo))|(?P<Years>\d+ ?years?)|(,?( and)? ?))+ (?P<Contents>.+?)$`)

type rmd struct {
	g *Bot
	meta TriggerMeta
}

func NewReminder(g *Bot) Trigger {
	var reminders []data.Reminder
	g.db.Model(&data.Reminder{}).Find(&reminders)
	rmd := rmd{
		g:g,
		meta: TriggerMeta{
			Disabled: false,
			Priority: 0,
			Name:     "reminder",
		},
	}
	for _, reminder := range reminders {
		go rmd.waitRemind(reminder)
	}
	return rmd
}

func (rmd rmd) waitRemind(r data.Reminder) {
	logrus.Debug("forked off reminder thread for %+v", r)
	<- time.After(time.Until(r.At))
	rmd.g.msgChan <- &irc.Message{
		Command: irc.PRIVMSG,
		Params:  []string{r.Location, r.Text},
	}
	rmd.g.db.Delete(&r)
}

func (rmd rmd) GetMeta() *TriggerMeta {
	return &rmd.meta
}

func (rmd rmd) Condition(g *Bot, msg *irc.Message) (shouldApply bool) {
	return msg.Command == irc.PRIVMSG && len(msg.Params) == 2 && reminderPattern.MatchString(msg.Params[1])
}

func (rmd rmd) Action(g *Bot, msg *irc.Message) (shouldContinue bool) {
	matches := getParams(msg.Params[1], reminderPattern)
	contents := ""
	duration := time.Nanosecond

	for k, v := range matches {
		if len(v) == 0 {
			continue
		}
		switch k {
		case "Seconds":
			i, _ := strconv.Atoi(strings.Split(v, " ")[0])
			duration += time.Duration(i) * time.Second
		case "Minutes":
			i, _ := strconv.Atoi(strings.Split(v, " ")[0])
			duration += time.Duration(i) * time.Minute
		case "Hours":
			i, _ := strconv.Atoi(strings.Split(v, " ")[0])
			duration += time.Duration(i) * time.Hour
		case "Days":
			i, _ := strconv.Atoi(strings.Split(v, " ")[0])
			duration += time.Duration(i) * time.Hour * 24
		case "Weeks":
			i, _ := strconv.Atoi(strings.Split(v, " ")[0])
			duration += time.Duration(i) * time.Hour * 24 * 7
		case "Months":
			i, _ := strconv.Atoi(strings.Split(v, " ")[0])
			duration += time.Duration(i) * time.Hour * 24 * 30
		case "Years":
			i, _ := strconv.Atoi(strings.Split(v, " ")[0])
			duration += time.Duration(i) * time.Hour * 24 * 365
		case "Contents":
			contents = v
		}
	}

	// remind me in 1 second to do the wash
	contents = strings.TrimPrefix(contents, "to ")
	// remind me in 1 second that the wash needs done
	contents = strings.TrimPrefix(contents, "that ")
	// remind me in 1 second that I need to do the wash
	// -> You need to do the wash
	if strings.HasPrefix(contents, "I ") {
		contents = strings.Replace(contents, "I ", "You", 1)
	}

	contents = fmt.Sprintf("%s: %s", msg.Name, contents)

	r := data.Reminder{
		Location: mirrorMsg(g, msg),
		Text: contents,
		At: time.Now().Add(duration),
	}
	go rmd.waitRemind(r)
	rmd.g.db.Create(&r)

	g.msgChan <- &irc.Message{
		Command: irc.PRIVMSG,
		Params:  []string{mirrorMsg(g, msg), fmt.Sprintf("I'll send you the reminder '%s' at %s", r.Text, r.At.Format("2006-01-02 15:04:05 -0700 MST"))},
	}
	return false
}


func getParams(haystack string, pattern *regexp.Regexp) (paramsMap map[string]string) {

	match := pattern.FindStringSubmatch(haystack)

	paramsMap = make(map[string]string)
	for i, name := range pattern.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}