package runtimevars

import (
	"os"
	"strings"
)

type IRV interface {
	Load()
	Clear()
	Get(key string) string
	Add(key, val string) string
}

type RV struct {
	Map map[string]string
}

func (rv *RV) Load() {
	runVars := make(map[string]string)
	for _, ev := range os.Environ() {
		pair := strings.SplitN(ev, "=", 2)
		runVars[pair[0]] = pair[1]
	}
	rv.Map = runVars
}

func (rv *RV) Get(key string) string {
	val := rv.Map[key]
	return val
}

func (rv *RV) Clear() {
	for key := range rv.Map {
		delete(rv.Map, key)
	}

}

func (rv *RV) Add(key, val string) {
	if rv.Map == nil {
		rv.Map = make(map[string]string)
	}
	rv.Map[key] = val
}
