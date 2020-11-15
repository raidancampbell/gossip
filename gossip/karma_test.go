package gossip

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/sorcix/irc.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestKarmaCounterCond(t *testing.T) {
	m := irc.Message{
		Prefix:  nil,
		Command: irc.PRIVMSG,
		Params:  []string{"#channel", "C++"},
	}
	assert.True(t, karmaCounter.Cond(nil, &m))

	m.Params[1] = "C--"
	assert.True(t, karmaCounter.Cond(nil, &m))

	m.Params[1] = "C"
	assert.False(t, karmaCounter.Cond(nil, &m))

	m.Params[1] = "middle of the sentence notepad++ is dropped"
	assert.True(t, karmaCounter.Cond(nil, &m))

	m.Params[1] = "middle of the sentence notepad-- is dropped"
	assert.True(t, karmaCounter.Cond(nil, &m))

	m.Params[1] = "middle of the (multi word)++ is dropped"
	assert.True(t, karmaCounter.Cond(nil, &m))

	m.Params[1] = "middle of the (multi word)-- is dropped"
	assert.True(t, karmaCounter.Cond(nil, &m))

	m.Params[1] = "ab+-+-+-+-+-+-a"
	assert.False(t, karmaCounter.Cond(nil, &m))
}

func testSetInitialKarma(mock sqlmock.Sqlmock, val int, object, location string) {
	tm, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", "2020-11-11 21:10:45.708593-07:00")
	if err != nil {
		panic(err)
	}
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "object", "value", "location"}).
		AddRow("1", tm, tm, tm, object, val, location)

	mock.
		ExpectQuery("SELECT \\* FROM `karmas` WHERE `karmas`.`object` = \\? AND `karmas`.`location` = \\? AND `karmas`.`deleted_at` IS NULL ORDER BY `karmas`.`id` LIMIT 1").
		WithArgs(object,location).
		WillReturnRows(rows)
}

func testSetUpdatedKarma(mock sqlmock.Sqlmock, val int, object, location string) {
	tm, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", "2020-11-11 21:10:45.708593-07:00")
	if err != nil {
		panic(err)
	}

	mock.ExpectExec("UPDATE `karmas` SET `created_at`=\\?,`updated_at`=\\?,`deleted_at`=\\?,`object`=\\?,`value`=\\?,`location`=\\? WHERE `id` = \\?").
		WithArgs(tm, sqlmock.AnyArg(), tm, object, val, location, 1).WillReturnResult(sqlmock.NewResult(0, 1))
}

func TestKarmaCounterAct(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer assert.Nil(t, mock.ExpectationsWereMet())

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
		SkipInitializeWithVersion: true,
	}), nil)
	if err != nil {
		panic(err)
	}
	b := &Bot{db:db}

	m := irc.Message{
		Prefix:  nil,
		Command: irc.PRIVMSG,
		Params:  []string{"#channel", "C++"},
	}

	// C will have karma of 2 in #channel
	testSetInitialKarma(mock, 2, "C", "#channel")
	// after someone says C++, expected karma is 3
	testSetUpdatedKarma(mock, 3, "C", "#channel")
	assert.True(t, karmaCounter.Act(b, &m))

	m.Params[1] = "middle of the (multi word)++ is used"
	testSetInitialKarma(mock, 0, "multi word", "#channel")
	testSetUpdatedKarma(mock, 1, "multi word", "#channel")
	assert.True(t, karmaCounter.Act(b, &m))

	m.Params[1] = "overly complex (we're talking about sentences here (english sentences)++ because they are the best) sentences are parsed naively"
	testSetInitialKarma(mock, 0, "english sentences", "#channel")
	testSetUpdatedKarma(mock, 1, "english sentences", "#channel")
	assert.True(t, karmaCounter.Act(b, &m))
}