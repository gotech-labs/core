package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type frame struct {
	name     string
	file     string
	function string
	line     int
}

type stackFrames []frame

func (st *stackFrames) traceLines() []string {
	var lines = make([]string, 0, len(*st)-1)
	for _, s := range *st {
		lines = append(lines, fmt.Sprintf("    at %s:%d (%s)", s.file, s.line, s.function))
	}
	return lines
}

type stack []uintptr

func (s *stack) StackFrames() stackFrames {
	frames := make([]frame, len(*s))
	for i, pc := range *s {
		fn := runtime.FuncForPC(pc)

		// fn.Name() is like one of these:
		// - "github.com/xxx/package.FuncName"
		// - "github.com/xxx/package.Receiver.MethodName"
		// - "github.com/xxx/package.(*PtrReceiver).MethodName"
		name := fn.Name()
		file, line := fn.FileLine(pc - 1)
		withoutPath := name
		if pos := strings.LastIndex(name, "/"); pos > 0 {
			withoutPath = name[pos+1:]
		}
		withoutPackage := withoutPath
		if pos := strings.Index(withoutPath, "."); pos > 0 {
			withoutPackage = withoutPath[pos+1:]
		}
		function := withoutPackage
		for _, target := range []string{"(", "*", ")"} {
			function = strings.Replace(function, target, "", 1)
		}
		frames[i] = frame{
			name:     name,
			file:     file,
			function: function,
			line:     line,
		}
	}
	return frames
}

func callers(skip int) *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0 : n-2]
	return &st
}
