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

package cmd

import (
	"github.com/faas-facts/fact/fact"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"runtime"
	"strconv"
)

var collect = &cobra.Command{
	Use:   "tcp [port] [max-connections] (threads)",
	Short: "starts tcp collection mode",
	Long:  `starts a tcp based collection server`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			fail(err, "port malformed")
		}

		connections, err := strconv.Atoi(args[1])
		if err != nil {
			fail(err, "max-connections malformed")
		}

		var threads int

		if len(args) > 2 {
			threads, err = strconv.Atoi(args[2])
			if err != nil {
				fail(err, "thread conf malformed")
			}
		} else {
			threads = runtime.NumCPU()
		}

		log.Infof("starting tcp on \":%d\" with %d max-connections and %d threads", port, connections, threads)

		collector := fact.NewTCPCollector(port, threads, connections)

		if viper.GetBool("continues") {
			registerObserver(collector.ResultCollector)
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for _ = range c {
				collector.Close()
			}
		}()

		collector.Listen()
		log.Info("closed tcp server trigger final write")
		write(collector.ResultCollector)

	},
}
