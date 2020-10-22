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
