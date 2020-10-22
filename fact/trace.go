package fact

import (
	"fmt"
)

type Platform string

const (
	AWS Platform = "AWS" //AWS Lambda
	ICF          = "ICF" // IBM Cloud Functions
	GCF          = "GCF" // Google Cloud Functions
	ACF          = "ACF" // Microsoft Azure Cloud Functions
	OW           = "OW"  // OpenWhisk
)

//merge a partial trace with the same TraceID to this trace. Older timestamps have precidence.
func (t *Trace) Merge(partial *Trace) error {
	if t.ID != partial.ID {
		return fmt.Errorf("partial trace dose not match")
	}

	if t.Status == 0 {
		t.Status = partial.Status
	}

	if t.Platform == "" {
		t.Platform = partial.Platform
	}

	if t.Region == "" {
		t.Region = partial.Region
	}

	if t.Memory == 0 {
		t.Memory = partial.Memory
	}

	if t.Runtime == "" {
		t.Runtime = partial.Runtime
	}

	if t.Env == nil {
		t.Env = partial.Env
	} else if partial.Env != nil {
		for k, v := range partial.Env {
			if _, ok := t.Env[k]; !ok {
				t.Env[k] = v
			}
		}
	}

	if t.ContainerID == "" {
		t.ContainerID = partial.ContainerID
	}

	if t.HostID == "" {
		t.HostID = partial.HostID
	}

	if t.BootTime == nil {
		t.BootTime = partial.BootTime
	}

	if t.StartTime == nil {
		t.StartTime = partial.StartTime
	}

	if t.EndTime == nil {
		t.EndTime = partial.EndTime
	}

	if t.ExecutionLatency == nil {
		t.ExecutionLatency = partial.ExecutionLatency
	}

	//XXX: we merge cost by using the oldest mesured cost ;)
	if t.Cost == 0 {
		t.Cost = partial.Cost
	} else {
		if t.Timestamp.AsTime().Before(partial.Timestamp.AsTime()) {
			t.Cost = partial.Cost
		}
	}

	if t.Tags == nil {
		t.Tags = partial.Tags
	} else if partial.Tags != nil {
		for k, v := range partial.Tags {
			if _, ok := t.Tags[k]; !ok {
				t.Tags[k] = v
			}
		}
	}

	if t.Logs == nil {
		t.Logs = partial.Logs
	} else if partial.Logs != nil {
		for k, v := range partial.Logs {
			if _, ok := t.Logs[k]; !ok {
				t.Logs[k] = v
			}
		}
	}

	if t.Args == nil {
		t.Args = partial.Args
	} else if partial.Args != nil {
		t.Args = append(t.Args, partial.Args...)
	}

	return nil
}
