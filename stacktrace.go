package gcplog

import (
	"fmt"
	"runtime"
	"strings"
)

type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

type StackTrace []Frame

func (st *StackTrace) String() string {
	var out []string
	for _, frame := range *st {
		out = append(out, fmt.Sprintf("at %s (%s:%d)", frame.name(), frame.file(), frame.line()))
	}
	return strings.Join(out, "\n")
}

const (
	skipCaller  = 4
	callerDepth = 32
)

//Callers returns the stacktrace
func Callers() *StackTrace {
	var pcs [callerDepth]uintptr
	n := runtime.Callers(skipCaller, pcs[:])
	st := pcs[0:n]

	var f StackTrace = make([]Frame, len(st))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((st)[i])
	}
	return &f
}
