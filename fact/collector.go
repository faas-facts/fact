package fact

import (
	"io"
	"io/ioutil"
	"sync"

	"github.com/golang/protobuf/proto"
)

type TraceObserver interface {
	Observe(trace *Trace) //
	Close()
}

type ResultCollector struct {
	sync.RWMutex
	traces    []*Trace
	observers []TraceObserver
	updates   chan *Trace
}

func NewCollector() *ResultCollector {
	rc := &ResultCollector{
		traces:    make([]*Trace, 0),
		observers: make([]TraceObserver, 0),
		updates:   make(chan *Trace),
	}

	go rc.startObservers()

	return rc
}

func (c *ResultCollector) AddObserver(observer TraceObserver) {
	c.observers = append(c.observers, observer)
}

func (c *ResultCollector) Decode(reader io.Reader) error {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var t Trace
	err = proto.Unmarshal(buf, &t)
	if err != nil {
		//XXX: what to do!?
		panic(err)
	}

	c.Add(&t)

	return nil
}

func (c *ResultCollector) Add(t *Trace) {
	c.Lock()
	c.traces = append(c.traces, t)
	c.Unlock()

	c.updates <- t
}

func (c *ResultCollector) merge(in []*Trace) []*Trace {
	traces := make(map[string]*Trace)
	for _, t := range in {
		if vv, ok := traces[t.ID]; ok {
			//XXX: we strongly assume that the IDs of both traces are equal ;)
			_ = vv.Merge(t)
		} else {
			traces[t.ID] = t
		}
	}

	result := make([]*Trace, 0)
	for _, v := range traces {
		result = append(result, v)
	}

	return result
}

//Write merges and writes all collected traces to the provider provider. Warning! This action will delete all collected traces after successful write operations.
func (c *ResultCollector) Write(writer TraceWriter) error {
	c.Lock()
	defer c.Unlock()
	traces := c.merge(c.traces)

	err := writer.Write(traces)
	if err != nil {
		return err
	}
	//if write is successful we delete the old traces to free up memory
	c.traces = make([]*Trace, 0)
	return nil
}

func (c *ResultCollector) GetTraces() []*Trace {
	c.Lock()
	defer c.Unlock()
	dump := make([]*Trace, len(c.traces))
	copy(dump, c.traces)
	return dump
}

func (c *ResultCollector) startObservers() {
	for t := range c.updates {
		for _, observer := range c.observers {
			observer.Observe(t)
		}
	}
}

func (c *ResultCollector) Close() {
	close(c.updates)
	for _, observer := range c.observers {
		observer.Close()
	}

}
