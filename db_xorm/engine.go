package db_xorm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"os"
	"xorm.io/xorm"
	xormLog "xorm.io/xorm/log"
)

var DB *xorm.Engine

func Init() error {
	var err error
	// XORM创建引擎
	DB, err = xorm.NewEngine("mysql", "root:root@(127.0.0.1:3308)/media_prime?charset=utf8mb4&interpolateParams=true")
	if err != nil {
		return err
	}
	// 创建自定义的日志实例
	log := logrus.New()
	log.Out = os.Stdout
	// 将日志实例设置到XORM的引擎中
	DB.SetLogger(&TracerLogger{
		logger:  log,
		level:   xormLog.LOG_DEBUG,
		showSQL: true,
	})
	return nil
}

type TracerLogger struct {
	logger  *logrus.Logger
	level   xormLog.LogLevel
	showSQL bool
	span    opentracing.Span
}

func (l *TracerLogger) BeforeSQL(ctx xormLog.LogContext) {
	l.span, _ = opentracing.StartSpanFromContext(ctx.Ctx, "xorm")
}

func (l *TracerLogger) AfterSQL(ctx xormLog.LogContext) {
	defer l.span.Finish()
	var sessionPart string
	l.span.LogFields(opentracingLog.String("db.statement", ctx.SQL))
	l.span.LogFields(opentracingLog.Object("db.args", ctx.Args))
	l.span.LogFields(opentracingLog.String("db.type", "sql"))
	l.span.LogFields(opentracingLog.Object("db.execute_time", ctx.ExecuteTime))
	if ctx.ExecuteTime > 0 {
		l.logger.Infof("[SQL]%s %s %v - %v", sessionPart, ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.logger.Infof("[SQL]%s %s %v", sessionPart, ctx.SQL, ctx.Args)
	}
}

func (l *TracerLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
	return
}

func (l *TracerLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
	return
}

func (l *TracerLogger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
	return
}

func (l *TracerLogger) Warnf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
	return
}

func (l *TracerLogger) Level() xormLog.LogLevel {
	return l.level
}

func (l *TracerLogger) SetLevel(lv xormLog.LogLevel) {
	l.level = lv
	return
}

func (l *TracerLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (l *TracerLogger) IsShowSQL() bool {
	return l.showSQL
}
