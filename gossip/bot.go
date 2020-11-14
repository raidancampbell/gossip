package gossip

import (
	"fmt"
	"github.com/raidancampbell/gossip/conf"
	"github.com/raidancampbell/gossip/data"
	"github.com/raidancampbell/libraidan/pkg/rruntime"
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

type Bot struct {
	addr string
	nick string
	channels []string
	msgChan chan *irc.Message
	c *irc.Conn
	joinChannels *sync.Once // todo: move this into its own struct
	triggers []Trigger
	db *gorm.DB
}

func New(n conf.Network, nick string) *Bot {
	b := &Bot{
		addr: fmt.Sprintf("%s:%d", n.Host, n.Port),
		nick: nick,
		channels: n.Channels,
		msgChan: make(chan *irc.Message),
		joinChannels: &sync.Once{},
		triggers: []Trigger{pingPong, joinChans, invite, userPingPong, htmlTitle, quit, part, rename, karmaCounter, KarmaBest, KarmaWorst},
	}
	/* Feature todo:
	[X] control its nick
	[ ] owner authorization
	[X] reminders (needs state)
	[X] karma (needs state)
	[X] source
	[ ] wolfram
	[ ] youtube
	[X] part
	[ ] feature toggles (needs interface/impl for features)
	*/
	var err error
	b.db, err = gorm.Open(sqlite.Open("gossip.db"), &gorm.Config{
		//TODO: not really doing anything
		//Logger: &data.GormLogger{LogLevel: logger.Info},
	})
	b.db = b.db.Debug()
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	b.db.AutoMigrate(&data.Reminder{})
	b.db.AutoMigrate(&data.Karma{})
	b.triggers = append(b.triggers, NewReminder(b))

	return b
}

// will block
func (g *Bot) Begin() {
	c, err := irc.Dial(g.addr)
	if err != nil {
		logrus.WithError(err).Errorf("unable to dial IRC addr '%s'", g.addr)
	}
	logrus.Infof("successfully connected to network '%s'", g.addr)
	g.c = c

	go g.encodeLoop()
	go func() {
		g.msgChan <- &irc.Message{
			Command: irc.USER,
			Params:  []string{g.nick, "0", "*", fmt.Sprintf("%s-irc-bot", g.nick)},
		}
		g.msgChan <- &irc.Message{
			Command: irc.NICK,
			Params:  []string{g.nick},
		}
	}()

	g.decodeLoop()
}

func (g *Bot) encodeLoop() {
	for msg := range g.msgChan {
		if g.c == nil {
			logrus.Infof("connection closed, exiting %s", rruntime.GetMyFuncName())
			return
		}
		logrus.WithFields(logrus.Fields{"message": fmt.Sprintf("%+v", msg), "raw_message": msg}).Debug("outgoing message")
		err := g.c.Encode(msg)
		if err != nil {
			logrus.WithError(err).Error("error during message encoding")
		}
	}
}


func (g *Bot) decodeLoop() {
	for {
		if g.c == nil {
			logrus.Infof("connection closed, exiting %s...", rruntime.GetMyFuncName())
			return
		}
		msg, err := g.c.Decode()
		if err != nil {
			logrus.WithError(err).Error("error during message decoding")
		}
		if msg == nil {
			logrus.Infof("no message to decode. exiting...")
			return
		}
		logrus.WithFields(logrus.Fields{"message": fmt.Sprintf("%+v", msg), "raw_message": msg}).Debug("incoming message")
		// each incoming message gets its own goroutine
		// just in case I really screw up and something hangs/dies,
		// so that pingpong still lives on
		go g.handleTriggers(msg)
	}
}

func (g *Bot) handleTriggers(msg *irc.Message) {
	for _, trigger := range g.triggers {
		if !trigger.GetMeta().Disabled && trigger.Condition(g, msg) {
			shouldContinue := trigger.Action(g, msg)
			if !shouldContinue {
				break
			}
		}
	}
}