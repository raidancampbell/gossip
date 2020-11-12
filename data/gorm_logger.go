package data

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct{
	LogLevel logger.LogLevel
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *g
	newlogger.LogLevel = level
	return &newlogger
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	logrus.WithField("caller", "gorm").Info(ctx, s, i)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	logrus.WithField("caller", "gorm").Warn(ctx, s, i)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	logrus.WithField("caller", "gorm").Error(ctx, s, i)
}

// unimplemented, I don't really care
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	return
}
