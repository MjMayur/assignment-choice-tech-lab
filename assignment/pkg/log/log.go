package log

import (
	"fmt"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func init() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// Removed everything before "falcon_be/"
		if idx := strings.Index(file, "falcon_be/"); idx != -1 {
			file = file[idx+len("falcon_be/"):] // cut falcon_be/
		}
		return file + ":" + strconv.Itoa(line)
	}
}

// -------------------- Logging Wrappers --------------------

func Info(msg string, requestID string) {
	zlog.Info().Str("RequestID", requestID).Caller(1).Msg(msg)
}

func InfoWithData(msg string, data interface{}, requestID string) {
	zlog.Info().Str("RequestID", requestID).Caller(1).
		Msg(msg + fmt.Sprintf(" %+v", data))
}

func Error(msg string, requestID string) {
	zlog.Error().Str("RequestID", requestID).Caller(1).Msg(msg) //usage : logging.Error("msg", "reqID")
}

func ErrorWithData(msg string, data interface{}, requestID string) {
	zlog.Error().Str("RequestID", requestID).Caller(1).
		Msg(msg + fmt.Sprintf(" %+v", data))
}

func Debug(msg string, requestID string) {
	zlog.Debug().Str("RequestID", requestID).Caller(1).Msg(msg)
}

func DebugWithData(msg string, data interface{}, requestID string) {
	zlog.Debug().Str("RequestID", requestID).Caller(1).
		Msg(msg + fmt.Sprintf(" %+v", data))
}

func Print(v ...interface{}) {
	zlog.Print(v...)
}

func Printf(format string, v ...interface{}) {
	zlog.Printf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	zlog.Debug().Msgf(format, v...)
}
