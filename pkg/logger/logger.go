package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BlockChain-Passion/go-project/pkg/runtimevars"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ILogger interface {
	Setup(rv runtimevars.IRV) error
	GetLogrous() *logrus.Logger
	SetLevel(lvl string) error
}

type DefaultLogger struct {
	Logg      *logrus.Logger
	LogLvlKey string `validate:"required"`
}

func (dl *DefaultLogger) Setup(rv runtimevars.IRV) error {
	// validate the struct
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(dl); err != nil {
		return err
	}
	// check rv
	// if rv == nil {
	// 	rv.Load()
	// }

	//setup Logger
	dl.Logg = logrus.New()
	dl.Logg.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	dl.Logg.SetReportCaller(true)

	dl.Logg.SetOutput(os.Stdout)

	err1 := dl.SetLogrous(dl.LogLvlKey)

	if err1 != nil {
		dl.Logg.Info("using default log level ", err1)
		return err1
	}

	// logLvl, err := logrus.ParseLevel(rv.Get(dl.LogLvlKey))

	// if err != nil {
	// 	dl.Logg.Info("using default log Level")
	// 	logLvl = 4
	// }
	// dl.Logg.SetLevel(logLvl)
	return nil

}

func (dl *DefaultLogger) GetLogrous() *logrus.Logger {
	return dl.Logg
}

func (dl *DefaultLogger) SetLogrous(lvl string) error {
	logLvl, err := logrus.ParseLevel(lvl)
	if err != nil {
		dl.Logg.Info("error while getting ParseLevel ", err)
		logLvl = 4
		return err
	}
	dl.Logg.SetLevel(logLvl)
	return nil
}

// Structured Logger is simple and powerfull implementation of a custom structed
// logger backed on logrous . Designed for context-based http router.

type StructuredLogger struct {
	Logger *logrus.Logger
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed_ms":   float64(elapsed.Nanoseconds()) / 1000000.0,
	})
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	}

	logFields["herokou_api_id"] = os.Getenv("HEROKU_APP_ID")
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method
	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()
	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)
	entry.Logger.Infoln("request started")

	return entry
}

// Helper methods used by the applications to get the request-scoped logger entry and set additional fields between handlers.
// This is a useful pattern to use to set state on the entry as it passes through the handlers chain , which at any point can be logged with a
// call to .Print(), Info(). etc.

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r)
	if entry != nil {
		return entry.(*StructuredLoggerEntry).Logger
	}
	return StructuredLoggerEntry{Logger: logrus.NewEntry(logrus.StandardLogger())}.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
