package fact_test

import (
	"fmt"
	"github.com/faas-facts/fact/fact"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestTCPCollector_Basic(t *testing.T) {

	traces := GenTraces(20)

	port := 9999
	workers := 3
	collector := fact.NewTCPCollector(port, workers, 5)

	go collector.Listen()

	sliceLenght := len(traces) / workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go tcpClient(t, &wg, port, traces[sliceLenght*i:sliceLenght*(i+1)])
	}

	wg.Wait()

	reflect.DeepEqual(traces, collector.GetTraces())

}

func tcpClient(t *testing.T, wg *sync.WaitGroup, collectorPort int, traces []*fact.Trace) {

	for _, trace := range traces {
		t.Logf("sending %+v", trace.ID)
		<-time.After(time.Duration(rand.Intn(300)) * time.Millisecond)
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", collectorPort))
		check(t, err)

		data, err := proto.Marshal(trace)
		check(t, err)

		_, err = conn.Write(data)
		check(t, err)

		err = conn.Close()
		check(t, err)

	}
	wg.Done()
}
