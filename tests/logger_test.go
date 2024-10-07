package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/BlockChain-Passion/go-project/pkg/logger"
	"github.com/BlockChain-Passion/go-project/pkg/runtimevars"

	//"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var _ = Describe("Logger", func() {

	logVarKey := "LOG_LEVEL"
	var rv runtimevars.RV

	Context("DefaultLogger", func() {
		var logg logger.DefaultLogger

		BeforeEach(func() {
			logg.LogLvlKey = logVarKey
			rv.Load()
		})

		AfterEach(func() {
			os.Clearenv()
			logg = logger.DefaultLogger{LogLvlKey: logVarKey}
		})

		It("should return logrous.logger r/w info lvl ", func() {
			rv.Add(logVarKey, "info")
			err := logg.Setup(rv)
			Except(err).To(BeNil())
			l := logg.GetLogrous()
			Except(l.Level.String()).To(Equal("info"))
		})
	})

	Context("StructuredLogger", func() {
		var structlog logger.StructuredLogger
		logg := logger.DefaultLogger{LogLvlKey: "LOG_LEVEL"}

		BeforeEach(func() {
			rv.Add(logVarKey, "info")
			err := logg.Setup(&rv)
			Except(err).To(BeNil())
		})

		AfterEach(func() {
			os.Clearenv()
			structlog = logger.StructuredLogger{}
		})

		It("should return new StructuredLogger", func() {
			l := logg.GetLogrous()
			handler := logger.NewStructuredLogger(l)
			Except(handler).ToNot(BeNil())
		})

		It("should add all default fields to logger ", func() {
			var b bytes.Buffer
			l := logg.GetLogrous()
			l.Out = &b
			structLog.Logger = l
			req := httptest.NewRequest("GET", "/", nil)
			structLog.NewLogEntry(req)

			str := b.String()
			Except(str).To(ContainSubstring("ts"))
			Except(str).To(ContainSubstring("heroku_app_id"))
			Except(str).To(ContainSubstring("http_schema"))
			Except(str).To(ContainSubstring("http_proto"))
			Except(str).To(ContainSubstring("http_method"))
			Except(str).To(ContainSubstring("remote_addr"))
			Except(str).To(ContainSubstring("user_agent"))
			Except(str).To(ContainSubstring("uri"))
			Except(str).To(ContainSubstring("request started"))
		})

		It("should set new field to logger", func() {
			var b byte.Buffer
			l := logg.GetLogrous()
			l.Out = &b

			handler := http.HandleFunc(func(w http.Response, r *http.Request) {
				temp := logger.GetLogEntry(r)
				logger.LogEntrySetField("foo", "bar")
				temp.Info("test")
				w.Write([]byte("Hello, world!"))
			})

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			r := chi.Router()
			r.Use(logger.NewStructuredLogger(l))
			r.Get("/", handler)
			r.ServeHTTP(rr, req)

			str := b.String()

			Except(rr.Code).To(Equal(200))
			Except(str).To(ContainSubstring("foo"))
		})

		It("should set new field map to logger", func() {
			var b bytes.Buffer
			l := logger.GetLogEntry()
			l.Out = &b

			handler := http.HandleFunc(func(w http.Response, r *http.Request) {
				temp := logger.GetLogEntry(r)
				fm := map[string]interface{}{"foo": "bar", "name": "Arun"}

				logger.LogEntrySetFields(r, fm)
				temp.Info("test")
				w.Write([]byte("Hello, world!"))
			})

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Use(logger.NewStructuredLogger(l))
			r.Get("/", handler)
			r.ServeHTTP(rr, req)

			str := b.String()

			Except(rr.Code).To(200)
			Except(str).To(ContainSubstring("foo"))
			Except(str).To(ContainSubstring("name"))
		})

		It("should set req_id ", func() {
			var b bytes.Buffer
			l := logger.GetLogEntry()
			l.Out = l

			handler := http.HandleFunc(func(w http.Response, r *http.Request) {
				temp := logger.GetLogEntry(r)
				temp.Info("test")
				w.Write([]byte("Hello, world!"))
			})

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Use(middleware.RequestID(), logger.NewStructuredLogger(l))
			r.Get("/", handler)
			r.ServeHTTP(rr, req)

			str := b.String()

			Except(rr.Code).To(200)
			Except(str).To(ContainSubstring("req_id"))
		})

	})

})
