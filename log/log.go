package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

// D 表示用于结构化日志记录的入门级数据的映射。
// type D map[string]interface{}
type D struct {
	Key   string
	Value interface{}
}

// Config log config.
type Config struct {
	Family string
	Host   string

	// stdout
	Stdout bool

	// file
	Dir string
	// buffer size
	FileBufferSize int64
	// MaxLogFile
	MaxLogFile int
	// RotateSize
	RotateSize int64

	// V Enable V-leveled logging at the specified level.
	V int32
	// Module=""
	// The syntax of the argument is a map of pattern=N,
	// where pattern is a literal file name (minus the ".go" suffix) or
	// "glob" pattern and N is a V level. For instance:
	// [module]
	//   "service" = 1
	//   "dao*" = 2
	// sets the V level to 2 in all Go files whose names begin "dao".
	Module map[string]int32
	// Filter tell log handler which field are sensitive message, use * instead.
	Filter []string
}

// Render render log output
type Render interface {
	Render(io.Writer, map[string]interface{}) error
	RenderString(map[string]interface{}) string
}

var (
	h Handler
	c *Config
)

var (
	_v        int
	_stdout   bool
	_dir      string
	_agentDSN string
	//_filter   logFilter
	//_module   = verboseModule{}
	_noagent bool
)

// Init create logger with context.
func Init(conf *Config) {
	//var isNil bool
	_dir = "/Users/liuhaigang/go/src/aha-api-server/logs/"
	if conf == nil {
		//isNil = true
		conf = &Config{
			Stdout: _stdout,
			Dir:    _dir,
			V:      int32(_v),
			Module: map[string]int32{"log_test": 1},
			//Filter: _filter,
			FileBufferSize: 1014,
			RotateSize:     1014,
			MaxLogFile:     100012,
		}
	}
	//if len(env.AppID) != 0 {
	//	conf.Family = env.AppID // for caster
	//}
	conf.Host = "localhost"
	if len(conf.Host) == 0 {
		host, _ := os.Hostname()
		conf.Host = host
	}
	var hs []Handler
	// when env is dev
	//if conf.Stdout || (isNil && (env.DeployEnv == "" || env.DeployEnv == env.DeployEnvDev)) || _noagent {
	//hs = append(hs, NewStdout())
	//}
	//if conf.Dir != "" {
	hs = append(hs, NewFile(conf.Dir, conf.FileBufferSize, conf.RotateSize, conf.MaxLogFile))
	//}
	//// when env is not dev
	//if !_noagent && (conf.Agent != nil || (isNil && env.DeployEnv != "" && env.DeployEnv != env.DeployEnvDev)) {
	//hs = append(hs, NewAgent(conf.Agent))
	//}
	h = newHandlers(conf.Filter, hs...)
	c = conf
}

// Info logs a message at the info log level.
func Info(format string, args ...interface{}) {
	h.Log(context.Background(), _infoLevel, KV(_log, fmt.Sprintf(format, args...)))
}

// Warn logs a message at the warning log level.
func Warn(format string, args ...interface{}) {
	h.Log(context.Background(), _warnLevel, KV(_log, fmt.Sprintf(format, args...)))
}

// Error logs a message at the error log level.
func Error(format string, args ...interface{}) {
	h.Log(context.Background(), _errorLevel, KV(_log, fmt.Sprintf(format, args...)))
}

// Infov logs a message at the info log level.
func Infov(ctx context.Context, args ...D) {
	h.Log(ctx, _infoLevel, args...)
}

// Warnv logs a message at the warning log level.
func Warnv(ctx context.Context, args ...D) {
	h.Log(ctx, _warnLevel, args...)
}

// Errorv logs a message at the error log level.
func Errorv(ctx context.Context, args ...D) {
	h.Log(ctx, _errorLevel, args...)
}

// KV return a log kv for logging field.
func KV(key string, value interface{}) D {
	return D{
		Key:   key,
		Value: value,
	}
}

func errIncr(lv Level, source string) {
	if lv == _errorLevel {
		//TODO 数据错误上报
		//errProm.Incr(source)
	}
}

// Close close resource.
func Close() (err error) {
	err = h.Close()
	h = _defaultStdout
	return
}
