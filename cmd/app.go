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
	"fmt"
	"github.com/faas-facts/fact/fact"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

var log *logrus.Entry

var (
	// Used for flags.
	cfgFile string
	output  string
	file    string

	writer fact.TraceWriter

	rootCmd = &cobra.Command{
		Use:              "fact",
		Short:            "Fact - FaaS Application & Component Tracer",
		Long:             `Fact - FaaS Application Component Tracer - is a tool to collect unified monitoring, logging, and tracing information from serverless applications.`,
		PersistentPreRun: CreateWriterOrFail,
	}
)

func CreateWriterOrFail(cmd *cobra.Command, args []string) {
	switch output {
	case "csv":
		writer = fact.NewCSVWriter()
	default:
		output = "csv"
		writer = fact.NewCSVWriter()

	}

	if writer == nil {
		log.Fatal("failed to create output")
	}

	if file == "" {
		file = fmt.Sprintf("./fact_%s.%s", time.Now().Format("2006_01_02_1504"), output)
	}

	_, err := os.Stat(file)
	var f *os.File
	var append bool
	if err != nil {
		append = false
		f, err = os.Create(file)
	} else {
		append = true
		f, err = os.Open(file)
	}

	if err != nil {
		log.Fatalf("could not open output %s - %+v", file, err)
	}

	writer.Open(f, append)

}

func Execute(logger *logrus.Entry) {
	log = logger

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fact.yaml)")

	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output File Format (required)")
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	rootCmd.MarkFlagRequired("output")

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Output File Path")
	viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))

	rootCmd.PersistentFlags().BoolP("continues", "c", false, "Writes to the Output continuously")
	viper.BindPFlag("continues", rootCmd.PersistentFlags().Lookup("continues"))

	rootCmd.AddCommand(collect)

}

type observer struct {
	bufferSize int
	traces     int
	collector  *fact.ResultCollector
}

func (o observer) Observe(trace *fact.Trace) {
	o.traces += 1
	if o.traces >= o.bufferSize {
		write(o.collector)
		o.traces = 0
	}
}

func write(collector *fact.ResultCollector) {
	err := collector.Write(writer)
	if err != nil {
		log.Errorf("failed to write traces %+v", err)

	}
}

func (o observer) Close() {
	write(o.collector)
}

func registerObserver(collector *fact.ResultCollector) {
	collector.AddObserver(&observer{
		bufferSize: 1024,
		traces:     32,
		collector:  collector,
	})
}

func er(msg interface{}) {
	log.Fatal(msg)
}

func fail(err error, msg interface{}) {
	log.Fatalf("failed %+v with %+v", msg, err)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".fact" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".fact")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
