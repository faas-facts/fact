package fact_test

import (
	"bytes"
	"github.com/faas-facts/fact/fact"
	"testing"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestCSVWriter_ToConsole(t *testing.T) {
	csv := fact.NewCSVWriter()

	var b bytes.Buffer
	csv.Open(&b)

	traces := GenSampleTraces()

	err := csv.Write(traces)
	check(t, err)

	t.Log(b.String())

}
