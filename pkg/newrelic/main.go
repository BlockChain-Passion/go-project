package newrelic

import (
	"fmt"
	"os"

	"github.com/BlockChain-Passion/go-project/pkg/runtimevars"
	"github.com/go-playground/validator/v10"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type INewrelic interface {
	Connect(rv runtimevars.IRV) error
	GetNewRelicaApp() *newrelic.Application
}

type DefaultNewRelic struct {
	AppNameKey     string `validate:"required"`
	LicenseKey     string `validate:"required"`
	DisplayNameKey string `validate:"required"`
	App            *newrelic.Application
}

func (nr *DefaultNewRelic) Connect(rv runtimevars.IRV) error {
	//validate struct
	validate := validator.New()
	if err := validate.Struct(nr); err != nil {
		return err
	}

	if rv == nil {
		rv.Load()
	}

	con, err := newrelic.NewApplication(
		newrelic.ConfigAppName(rv.Get(nr.AppNameKey)),
		newrelic.ConfigLicense(rv.Get(nr.LicenseKey)),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigDebugLogger(os.Stdout),
		func(cfg *newrelic.Config) {
			cfg.ErrorCollector.RecordPanics = true
			cfg.HostDisplayName = rv.Get(nr.DisplayNameKey)
		},
	)

	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
	nr.App = con

	return nil
}

func (nr *DefaultNewRelic) GetNewRelicaApp() *newrelic.Application {
	return nr.App
}

// func NewRelicWrapper(next http.Handler, newRelic INewRelic) http.Handler {
// 	return http.HandlerFunc(
// 		func(rw http.ResponseWriter, r *http.Request) {
// 			nrApp := newRelic.GetNewRelicaApp()
// 			if nrApp != nil {
// 				trx := nrApp.StartTransaction(r.Method + " " + r.RequestURI)
// 				defer trx.End()
// 				rw = txn.SetWebResponse(rw)
// 				r = newrelic.RequestWithTransactionContext(r, txn)
// 			}

// 			next.ServeHTTP(rw, r)
// 		}
// 	)
// }
