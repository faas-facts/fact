package fact

import (
	"encoding/csv"
	"io"
	"strconv"
)

type CSVWriter struct {
	sink io.Writer
}

func (c *CSVWriter) Open(writer io.Writer) {
	c.sink = writer
}

func NewCSVWriter() TraceWriter {
	return &CSVWriter{}
}

func (c *CSVWriter) Name() string {
	return "CSV"
}

func (c *CSVWriter) header() []string {
	return []string{
		"ID", "ChildOf", "Timestamp", "CId", "HId", "CStart", "ECost", "EStart", "ECode", "EEnd",
		"Provider", "Region", "COs", "CMem", "ELat",
	}
}

// CSV Writer will not output Logs or Args
func (c *CSVWriter) Write(traces []*Trace) error {
	header := c.header()
	length := len(header)

	envMap := make(map[string]int)
	for _, t := range traces {
		for k, _ := range t.Env {
			if _, ok := envMap[k]; !ok {
				envMap[k] = length
				header = append(header, "E_"+k)
				length += 1
			}
		}
	}

	tagMap := make(map[string]int)
	for _, t := range traces {
		for k, _ := range t.Tags {
			if _, ok := tagMap[k]; !ok {
				tagMap[k] = length
				header = append(header, "T_"+k)
				length += 1
			}
		}
	}

	w := csv.NewWriter(c.sink)

	err := w.Write(header)

	if err != nil {
		return err
	}

	for _, t := range traces {
		record := make([]string, length)
		//"ID","ChildOf","Timestamp","CId","HId","CStart","ECost","EStart","ECode","EEnd",
		//		"Provider","Region","COs","CMem","ELat",
		record[0] = t.ID
		record[1] = t.ChildOf
		record[2] = strconv.FormatInt(t.StartTime.GetSeconds(), 10)
		record[3] = t.ContainerID
		record[4] = t.HostID
		record[5] = strconv.FormatInt(t.BootTime.GetSeconds(), 10)
		record[6] = strconv.FormatFloat(float64(t.Cost), 'E', -1, 32)
		record[7] = strconv.FormatInt(t.StartTime.GetSeconds(), 10)
		record[8] = string(t.Status)
		record[9] = strconv.FormatInt(t.EndTime.GetSeconds(), 10)
		record[10] = t.Platform
		record[11] = t.Region
		record[12] = t.Runtime
		record[13] = string(t.Memory)
		record[14] = strconv.FormatInt(int64(t.ExecutionLatency.AsDuration()), 10)

		for k, v := range t.Env {
			record[envMap[k]] = v
		}

		for k, v := range t.Tags {
			record[tagMap[k]] = v
		}

		err := w.Write(record)
		if err != nil {
			//this is an issue??!
		}
	}
	w.Flush()

	return nil
}
