package gossip

import (
	"errors"
	"fmt"
	"github.com/raidancampbell/gossip/conf"
	"github.com/raidancampbell/gossip/data"
	"github.com/raidancampbell/libraidan/pkg/rruntime"
	"github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

type Bot struct {
	addr     string
	nick     string
	msgChan  chan *irc.Message
	c        *irc.Conn
	triggers []Trigger
	db       *gorm.DB
	cfg      *conf.Cfg
}

func New(cfg *conf.Cfg) *Bot {

	b := &Bot{
		addr:     fmt.Sprintf("%s:%d", cfg.Network.Host, cfg.Network.Port),
		nick:     cfg.Nick,
		msgChan:  make(chan *irc.Message),
		triggers: []Trigger{pingPong, onConnect, NewJoin(cfg.Network.Channels), invite, NewPush(cfg), userPingPong, htmlTitle, die, part, rename, karma, toggle, source},
		cfg:      cfg,
	}
	/* Feature todo:
	[X] control its nick
	[X] owner authorization
	[X] reminders (needs state)
	[X] karma (needs state)
	[X] source
	[ ] wolfram
	[ ] youtube
	[X] part
	[X] feature toggles (needs interface/impl for features)
	[X] persistent feature toggles
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
	b.db.AutoMigrate(&data.TriggerMeta{})
	b.triggers = append(b.triggers, NewReminder(b))

	// initialize trigger state
	for i := range b.triggers {
		m := data.TriggerMeta{
			Name: b.triggers[i].GetMeta().Name,
		}
		tx := b.db.Model(&m).Where(&m).First(&m)
		if tx.RowsAffected == 0 {
			b.db.Model(&m).Create(b.triggers[i].GetMeta())
			// read it back to get gorm's internal metadata/index
			b.db.Model(&m).Where(&m).First(&m)
		}
		b.triggers[i].Meta(&m)

	}
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

	for {
		err = g.decodeLoop()
		if !strings.Contains(err.Error(), "connection closed") || !strings.Contains(err.Error(), "connection timed out") {
			return
		}
		c, err := irc.Dial(g.addr)
		if err != nil {
			logrus.WithError(err).Errorf("unable to dial IRC addr '%s'", g.addr)
		}
		logrus.Infof("successfully connected to network '%s'", g.addr)
		g.c = c
	}
}

func (g *Bot) encodeLoop() {
	for msg := range g.msgChan {
		if g.c == nil {
			logrus.Infof("connection closed, exiting %s", rruntime.GetMyFuncName())
			return
		}
		if msg.Command != irc.PONG {
			logrus.WithFields(logrus.Fields{"message": fmt.Sprintf("%+v", msg), "raw_message": msg}).Debug("outgoing message")
		}
		err := g.c.Encode(msg)
		if err != nil {
			logrus.WithError(err).Fatal("error during message encoding")
		}
	}
}

func (g *Bot) decodeLoop() error {
	for {
		if g.c == nil {
			logrus.Infof("connection closed, exiting %s...", rruntime.GetMyFuncName())
			return errors.New("connection closed")
		}
		msg, err := g.c.Decode()
		if err != nil {
			logrus.WithError(err).Error("error during message decoding")
		}
		if msg == nil {
			logrus.Infof("no message to decode. exiting...")
			return fmt.Errorf("%w: nil message, exiting loop", err)
		}
		if msg.Command != irc.PING {
			logrus.WithFields(logrus.Fields{"message": fmt.Sprintf("%+v", msg), "raw_message": msg}).Debug("incoming message")
		}
		// each incoming message gets its own goroutine
		// just in case I really screw up and something hangs/dies,
		// so that pingpong still lives on
		go g.handleTriggers(msg)
	}
}

func (g *Bot) handleTriggers(msg *irc.Message) {
	for _, trigger := range g.triggers {
		if !trigger.GetMeta().Disabled && trigger.Condition(g, msg) {
			logrus.Tracef("matched on trigger '%s'", trigger.GetMeta().Name)
			shouldContinue := trigger.Action(g, msg)
			if !shouldContinue {
				break
			}
		}
	}
}
