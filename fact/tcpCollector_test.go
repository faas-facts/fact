/*
 *  MIT License
 *
 *  Copyright (c) 2021. Fact Contributors
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

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
