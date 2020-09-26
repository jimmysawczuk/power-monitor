package main

import (
	"os"
	"runtime"
	"time"

	"github.com/jimmysawczuk/tmpl/tmplfunc"
)

type goEnv struct {
	OS   string
	Arch string
	Ver  string
}

type payload struct {
	Hostname string
	GoEnv    goEnv

	now time.Time

	config struct {
		timestampAssets bool
	}
}

func newPayload(timestampAssets bool) (payload, error) {
	h, _ := os.Hostname()

	return payload{
		Hostname: h,
		GoEnv: goEnv{
			Ver:  runtime.Version(),
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		now: time.Now(),
		config: struct {
			timestampAssets bool
		}{
			timestampAssets: timestampAssets,
		},
	}, nil
}

func (o payload) tmplfuncs() map[string]interface{} {
	return map[string]interface{}{
		"asset": tmplfunc.AssetLoaderFunc(o.now, o.config.timestampAssets),
		"env":   tmplfunc.Env,

		"getJSON": tmplfunc.GetJSON,
		"jsonify": tmplfunc.JSONify,

		"now":        tmplfunc.NowFunc(o.now),
		"parseTime":  tmplfunc.ParseTime,
		"formatTime": tmplfunc.FormatTime,
		"timeIn":     tmplfunc.TimeIn,

		"safeHTML":     tmplfunc.SafeHTML,
		"safeHTMLAttr": tmplfunc.SafeAttr,
		"safeJS":       tmplfunc.SafeJS,
		"safeCSS":      tmplfunc.SafeCSS,

		"seq": tmplfunc.Seq,
		"add": tmplfunc.Add,
	}
}
