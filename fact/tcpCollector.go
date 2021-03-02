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

package fact

/*
 * MIT License
 *
 * Copyright (c) 2020 Sebastian Werner
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type TCPCollector struct {
	*ResultCollector
	pool *tcpPool
	port int
}

func NewTCPCollector(port, worker, maxConnections int) *TCPCollector {
	base := NewCollector()
	collector := &TCPCollector{
		ResultCollector: base,
		pool:            newPool(worker, maxConnections, base.Decode),
		port:            port,
	}

	return collector
}

func (t *TCPCollector) Listen() {
	t.pool.Listen(t.port)
}

func (t *TCPCollector) Close() {
	t.pool.Close()
	t.ResultCollector.Close()
}

type tcpPool struct {
	sync.Mutex
	workers        int
	maxConnections int
	closed         bool

	pendingConnections chan net.Conn
	done               chan struct{}
	decoder            func(io.Reader) error
	ln                 net.Listener
}

func newPool(w int, t int, decoder func(io.Reader) error) *tcpPool {
	return &tcpPool{
		workers:            w,
		maxConnections:     t,
		pendingConnections: make(chan net.Conn, t),
		done:               make(chan struct{}),
		decoder:            decoder,
	}
}

func (p *tcpPool) Close() {
	p.Lock()
	defer p.Unlock()

	p.closed = true
	close(p.done)
	close(p.pendingConnections)

	_ = p.ln.Close()
}

func (p *tcpPool) addTask(conn net.Conn) {
	p.Lock()
	if p.closed {
		p.Unlock()
		return
	}
	p.Unlock()

	p.pendingConnections <- conn
}

func (p *tcpPool) start() {
	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}
}

func (p *tcpPool) startWorker() {
	for {
		select {
		case <-p.done:
			return
		case conn := <-p.pendingConnections:
			if conn != nil {
				p.handleConn(conn)
				_ = conn.Close()
			}
		}
	}
}

func (p *tcpPool) Listen(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	p.ln = ln
	p.start()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v - killing listener", e)
			if !p.closed {
				p.Close()
			}
			return nil
		}

		p.addTask(conn)
	}

}

func (p *tcpPool) handleConn(conn net.Conn) {
	err := p.decoder(conn)
	if err != nil {
		log.Printf("failed to read connection %+v", conn)
	}
}
