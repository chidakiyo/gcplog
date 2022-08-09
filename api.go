package gcplog

import (
	"context"
	"sync"
)

var debug = false

func SetDebug() {
	debug = true
}

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

var logLevel = DEBUG // default

func SetLoglevel(level int) {
	logLevel = level
}

type logger struct {
	mu        sync.Mutex
	text      Writer
	json      Writer
	structure func(context.Context, string, interface{})
}

func (l *logger) outputTxt(ctx context.Context, severity, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.text(ctx, severity, format, a...)
}

func (l *logger) outputJson(ctx context.Context, severity, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !debug {
		l.json(ctx, severity, format, a...)
	} else {
		l.text(ctx, severity, format, a...)
	}
}

func (l *logger) outputStructure(ctx context.Context, severity string, target interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.structure(ctx, severity, target)
}

var lg = &logger{
	text:      outputText,
	json:      outputJSON,
	structure: outputStructure,
}

type text struct{}

var Text = text{}

type structure struct{}

var Structure = structure{}

// text log

func (text) Debugf(ctx context.Context, format string, a ...interface{}) {
	if logLevel > DEBUG {
		return
	}
	lg.outputTxt(ctx, "DEBUG", format, a...)
}

func (text) Infof(ctx context.Context, format string, a ...interface{}) {
	if logLevel > INFO {
		return
	}
	lg.outputTxt(ctx, "INFO", format, a...)
}

func (text) Warningf(ctx context.Context, format string, a ...interface{}) {
	if logLevel > WARNING {
		return
	}
	lg.outputTxt(ctx, "WARNING", format, a...)
}

func (text) Errorf(ctx context.Context, format string, a ...interface{}) {
	if logLevel > ERROR {
		return
	}
	lg.outputTxt(ctx, "ERROR", format, a...)
}

func (text) Criticalf(ctx context.Context, format string, a ...interface{}) {
	lg.outputTxt(ctx, "CRITICAL", format, a...)
}

// json log

func Debugf(ctx context.Context, format string, a ...interface{}) {
	if logLevel > DEBUG {
		return
	}
	lg.outputJson(ctx, "DEBUG", format, a...)
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	if logLevel > INFO {
		return
	}
	lg.outputJson(ctx, "INFO", format, a...)
}

func Warningf(ctx context.Context, format string, a ...interface{}) {
	if logLevel > DEBUG {
		return
	}
	lg.outputJson(ctx, "WARNING", format, a...)
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	if logLevel > ERROR {
		return
	}
	lg.outputJson(ctx, "ERROR", format, a...)
}

func Criticalf(ctx context.Context, format string, a ...interface{}) {
	if logLevel > CRITICAL {
		return
	}
	lg.outputJson(ctx, "CRITICAL", format, a...)
}

// Structured log

func (structure) Debugf(ctx context.Context, a interface{}) {
	if logLevel > DEBUG {
		return
	}
	lg.outputStructure(ctx, "DEBUG", a)
}

func (structure) Infof(ctx context.Context, a interface{}) {
	if logLevel > INFO {
		return
	}
	lg.outputStructure(ctx, "INFO", a)
}

func (structure) Warningf(ctx context.Context, a interface{}) {
	if logLevel > WARNING {
		return
	}
	lg.outputStructure(ctx, "WARNING", a)
}

func (structure) Errorf(ctx context.Context, a interface{}) {
	if logLevel > ERROR {
		return
	}
	lg.outputStructure(ctx, "ERROR", a)
}

func (structure) Criticalf(ctx context.Context, a interface{}) {
	if logLevel > CRITICAL {
		return
	}
	lg.outputStructure(ctx, "CRITICAL", a)
}
