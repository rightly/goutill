package logger

import (
	"github.com/rightly/goutill"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Log struct {
	*zap.SugaredLogger
	mu    *sync.Mutex
	path  string
	level string
	*rotate
	*archive
}

type rotate struct {
	*cron.Cron
	flag     bool
	format   string
	interval string
	size     string
}

type archive struct {
	*cron.Cron
	flag       bool
	dateBefore string
}

func New() *Log {
	const timeFormat = "060102-15"
	var (
		lc = time.FixedZone("UTC-8", 9*60*60)

		env = goutill.String.ToUpper(goutill.OS.ENV())
		sl  = &zap.SugaredLogger{}
	)

	switch env {
	case "LOCAL":
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil
		}
		sl = logger.Sugar()
	case "DEV":
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil
		}
		sl = logger.Sugar()
	case "STAGE":
		logger, err := zap.NewProduction()
		if err != nil {
			return nil
		}
		sl = logger.Sugar()
	case "PROD":
		logger, err := zap.NewProduction()
		if err != nil {
			return nil
		}
		sl = logger.Sugar()
	}

	return &Log{
		SugaredLogger: sl,
		rotate: &rotate{
			Cron:     cron.New(cron.WithLocation(lc)),
			flag:     true,
			format:   timeFormat,
			interval: "1h",
			size:     "10m",
		},
		archive: &archive{
			Cron:       cron.New(cron.WithLocation(lc)),
			flag:       true,
			dateBefore: "7d",
		},
		mu: &sync.Mutex{},
	}
}

func (l *Log) Write() {

}

func (r *rotate) setLogRotate() error {
	if !r.flag {
		return nil
	}

	return nil
}

func (a *archive) setLogArchive() error {
	if !a.flag {
		return nil
	}

	return nil
}
