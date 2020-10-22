package fact

import "io"

type TraceWriter interface {
	Name() string         //name of the writer, for ui purposes
	Open(io.Writer)       //attaches an io.writer used to write all traces to
	Write([]*Trace) error //writes the contens of a ResultCollector
}
