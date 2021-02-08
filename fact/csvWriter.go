package fact

import (
	"encoding/csv"
	"io"
	"strconv"
)

type CSVWriter struct {
	sink   io.Writer
	append bool
}

func (c *CSVWriter) Open(writer io.Writer, append bool) {
	c.sink = writer
	c.append = append
}

func NewCSVWriter() TraceWriter {
	return &CSVWriter{}
}

func (c *CSVWriter) Name() string {
	return "CSV"
}

func (c *CSVWriter) header() []string {
	return []string{
		"ID", "ChildOf", "Timestamp", "CId", "HId", "CStart",
		"ECost", "RStart", "EStart", "ECode", "EEnd", "REnd",
		"Version", "CVersion", "Provider", "Region",
		"COs", "CMem", "ELat", "RLat", "DLat", "TLat",
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
	if !c.append {
		err := w.Write(header)

		if err != nil {
			return err
		}
		c.append = true
	}

	for _, t := range traces {
		record := make([]string, length)
		//"ID", "ChildOf", "Timestamp",
		record[0] = t.ID
		record[1] = t.ChildOf
		record[2] = strconv.FormatInt(t.StartTime.GetSeconds(), 10)
		//"CId", "HId", "CStart",
		record[3] = t.ContainerID
		record[4] = t.HostID
		record[5] = strconv.FormatInt(t.BootTime.GetSeconds(), 10)
		//"ECost", "RStart","EStart",
		record[6] = strconv.FormatFloat(float64(t.Cost), 'E', -1, 32)
		record[7] = strconv.FormatInt(t.RequestStartTime.GetSeconds(), 10)
		record[8] = strconv.FormatInt(t.StartTime.GetSeconds(), 10)
		//"ECode", "EEnd","REnd",
		record[9] = string(t.Status)
		record[10] = strconv.FormatInt(t.EndTime.GetSeconds(), 10)
		record[11] = strconv.FormatInt(t.RequestEndTime.GetSeconds(), 10)
		//"Version","CVersion","Provider", "Region",
		record[12] = t.CodeVersion
		record[13] = t.ConfigVersion
		record[14] = t.Platform
		record[15] = t.Region
		//"COs", "CMem", "ELat","RLat","DLat","TLat",
		record[16] = t.Runtime
		record[17] = string(t.Memory)
		record[18] = strconv.FormatInt(int64(t.ExecutionLatency.AsDuration()), 10)
		record[19] = strconv.FormatInt(int64(t.RequestResponseLatency.AsDuration()), 10)
		record[20] = strconv.FormatInt(int64(t.ExecutionDelay.AsDuration()), 10)
		record[21] = strconv.FormatInt(int64(t.TransportDelay.AsDuration()), 10)

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
