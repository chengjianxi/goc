package log

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type HaoxinJSONFormatter struct {
	ServiceId string
}

type HaoxinLog struct {
	Timestamp string `json:"timestamp"`           // 时间戳
	Level     string `json:"level"`               // 日志级别，取以下四个值 DEBUG,INFO,WARN,ERROR，分别表示调试打印、正常日志、警告、错误。
	Message   string `json:"message"`             // 日志内容
	ServiceId string `json:"serviceid"`           // 服务 ID
	Duration  int64  `json:"duration,omitempty"`  // 持续时间，毫秒计数
	SID       string `json:"sid,omitempty"`       // 表示操作发起方,格式:”模块ID(2位)”+”年月日时分秒(14位”+”随机数(6位)”.示例:4020190409104630000101
	TID       string `json:"tid,omitempty"`       // 表示操作目的方,格式:”模块ID(2位)”+”年月日时分秒(14位”+”随机数(6位)”示例:4020190409104630000101
	UserId    string `json:"userid,omitempty"`    // 用户id
	MachineId string `json:"machineid,omitempty"` // 用户id
}

func (f *HaoxinJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	level := "INFO"
	switch entry.Level {
	case logrus.TraceLevel:
		fallthrough
	case logrus.DebugLevel:
		level = "DEBUG"
	case logrus.InfoLevel:
		level = "INFO"
	case logrus.WarnLevel:
		level = "WARN"
	case logrus.FatalLevel:
		level = "ERROR"
	case logrus.ErrorLevel:
		level = "ERROR"
	}

	// 持续时间
	duration, _ := entry.Data["duration"].(time.Duration)
	sid, _ := entry.Data["sid"].(string)
	tid, _ := entry.Data["tid"].(string)
	userid, _ := entry.Data["userid"].(string)
	machineid, _ := entry.Data["machineid"].(string)

	log := HaoxinLog{
		Timestamp: entry.Time.Format("2006-01-02 15:04:05.000"),
		Level:     level,
		Message:   entry.Message,
		ServiceId: f.ServiceId, // 服务 ID
		Duration:  int64(duration),
		SID:       sid,
		TID:       tid,
		UserId:    userid,
		MachineId: machineid,
	}

	// Note this doesn't include Time, Level and Message which are available on
	// the Entry. Consult `godoc` on information about those fields or read the
	// source of the official loggers.
	serialized, err := json.Marshal(log)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}
	return append(serialized, '\n'), nil
}

func InitLog(logPath string, svcName string) io.Writer {
	formatter := &HaoxinJSONFormatter{ServiceId: svcName}
	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.InfoLevel)
	// 当天日志名 MTA.log, 每隔 1 天轮转一个新文件
	//logpath := path.Join(c.Log.Path, "MTA.log")
	//logrus.SetOutput(NewRotateWriter(logpath))
	// 下面配置日志轮，每隔 1 天轮转一个新文件，保留最近 7 天的日志文件，多余的自动清理掉。
	// github.com/lestrrat-go/file-rotatelogs
	logpattern := path.Join(logPath, svcName+"-%Y-%m-%d.log")
	writer, _ := rotatelogs.New(
		logpattern,
		rotatelogs.WithLinkName(""),
		rotatelogs.WithMaxAge(time.Duration(7*24)*time.Hour),     // 设置文件清理前的最长保存时间，WithMaxAge 和 WithRotationCount二者只能设置一个
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour), // 设置日志分割的时间，隔多久分割一次
	)
	logrus.SetOutput(writer)

	return writer
}

func requestFields(request *http.Request) map[string]interface{} {
	return logrus.Fields{
		"sid":       request.Header.Get("sid"),
		"tid":       request.Header.Get("tid"),
		"userid":    request.Header.Get("userid"),
		"machineid": request.Header.Get("machineid"),
	}
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func InfoWithFields(fields map[string]interface{}, args ...interface{}) {
	logrus.WithFields(fields).Info(args...)
}

func InfoLogWithRequest(request *http.Request, args ...interface{}) {
	InfoWithFields(requestFields(request), args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func WarnWithFields(fields map[string]interface{}, args ...interface{}) {
	logrus.WithFields(fields).Warn(args...)
}

func WarnLogWithRequest(request *http.Request, args ...interface{}) {
	WarnWithFields(requestFields(request), args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func ErrorWithFields(fields map[string]interface{}, args ...interface{}) {
	logrus.WithFields(fields).Error(args...)
}

func ErrorWithRequest(request *http.Request, args ...interface{}) {
	ErrorWithFields(requestFields(request), args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func PanicWithFields(fields map[string]interface{}, args ...interface{}) {
	logrus.WithFields(fields).Panic(args...)
}

func PanicWithRequest(request *http.Request, args ...interface{}) {
	PanicWithFields(requestFields(request), args...)
}
