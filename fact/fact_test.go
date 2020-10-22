package fact_test

import (
	"encoding/base64"
	"github.com/faas-facts/fact/fact"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(23456789))

func genString() string {
	buf := make([]byte, 10)
	for i := 0; i < 10; i++ {
		buf[i] = byte(random.Int())
	}
	return base64.StdEncoding.EncodeToString(buf)
}

func genStringSlice(n int) []string {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names = append(names, genString())
	}
	return names
}

func genMapFromNames(names []string, n int) map[string]string {
	rmap := make(map[string]string)
	for i := 0; i < n; i++ {
		rmap[names[random.Intn(len(names))]] = genString()
	}

	return rmap
}

func genLogs(n int) map[uint64]string {
	log := make(map[uint64]string)
	for i := 0; i < n; i++ {
		log[uint64(random.Int())] = genString()
	}
	return log
}

//Generates n synthetic traces
func GenTraces(n int) []*fact.Trace {

	envs := genStringSlice(10)
	tags := genStringSlice(10)

	traces := make([]*fact.Trace, 0)

	for i := 0; i < n; i++ {
		startTime := time.Now().Add(time.Duration(random.Intn(100)) * time.Minute * -1)
		trace := GenTrace("", startTime, envs, tags)

		traces = append(traces, trace)
	}

	return traces
}

//Generates a simple synthetic traces with a fixed startTime, a permutation of env and tags set.
func GenTrace(childOf string, startTime time.Time, envs []string, tags []string) *fact.Trace {
	id, _ := uuid.NewRandom()

	EStart := startTime.Add(time.Duration(random.Intn(10)) * time.Minute)
	EEnd := EStart.Add(time.Duration(random.Intn(400)) * time.Millisecond)

	trace := &fact.Trace{
		ID:      id.String(),
		ChildOf: childOf,
		Timestamp: &timestamp.Timestamp{
			Seconds: startTime.Unix(),
			Nanos:   0,
		},
		ContainerID: genString(),
		HostID:      genString(),
		BootTime: &timestamp.Timestamp{
			Seconds: startTime.Unix(),
			Nanos:   0,
		},
		Cost: 0,
		StartTime: &timestamp.Timestamp{
			Seconds: EStart.Unix(),
			Nanos:   0,
		},
		Status: 100,
		EndTime: &timestamp.Timestamp{
			Seconds: EEnd.Unix(),
			Nanos:   0,
		},
		Platform: "AWS",
		Region:   "eu-west-1",
		Runtime:  "Linux 2.5.4.1 / Python 2.7",
		Memory:   256,
		ExecutionLatency: &duration.Duration{
			Seconds: EEnd.Sub(EStart).Milliseconds() / 1000,
		},
		Env:  genMapFromNames(envs, 4),
		Tags: genMapFromNames(tags, 2),
		Logs: genLogs(10),
		Args: genStringSlice(10),
	}
	return trace
}

//Generates a set of traces based on other traces, each new traces is a ChildOf one of the provided traces.
func GenTracesFromTraces(parantes []*fact.Trace, added int) []*fact.Trace {
	envs := genStringSlice(10)
	tags := genStringSlice(10)

	results := make([]*fact.Trace, 0)
	results = append(results, parantes...)

	for i := 0; i < added; i++ {
		parent := parantes[random.Intn(len(parantes))]
		start := parent.StartTime.AsTime().Add(10 * time.Minute)
		results = append(results, GenTrace(parent.ID, start, envs, tags))
	}

	return results
}

//Generates a full set of synthetic traces including child-of relationships
func GenSampleTraces() []*fact.Trace {
	traces := GenTraces(3)
	return GenTracesFromTraces(traces, 2)
}
